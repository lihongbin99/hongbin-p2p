import java.io.Closeable;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.ZoneOffset;
import java.time.ZonedDateTime;
import java.util.Date;
import java.util.Map;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.ConcurrentHashMap;

public class Server {

    public final static int SERVER_PORT           = 9090;
    public final static int NAT_TYPE_TEST_PORT_1  = 9091;
    public final static int NAT_TYPE_TEST_PORT_2  = 9092;
    public final static int SHOW_ADDRESS_AND_PORT = 9093;

    public final static Map<String, Socket> CACHE = new ConcurrentHashMap<>();
    public final static String SEPARATOR = ":";
    public final static Timer TIMER = new Timer();

    public static void main(String[] args) throws Exception {
        natTypeTest(NAT_TYPE_TEST_PORT_1);
        natTypeTest(NAT_TYPE_TEST_PORT_2);
        showAddressAndPort();

        ServerSocket serverSocket = new ServerSocket(SERVER_PORT);
        System.out.println("服务器启动成功");
        while (!Thread.interrupted()) {
            try {
                Socket socket = serverSocket.accept();
                if (CACHE.size() >= 20) {
                    close(socket);
                } else {
                    apply(socket);
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        close(serverSocket);
        TIMER.cancel();
    }

    private static void apply(Socket clientSocket) throws Exception {
        final int hashCode = clientSocket.hashCode();
        InputStream inputStream = clientSocket.getInputStream();

        String clientName = getString(inputStream);
        // 三分钟内没匹配成功就放弃了
        TIMER.schedule(new TimerTask() {
            @Override
            public void run() {
                synchronized (CACHE) {
                    Socket socket = CACHE.get(clientName);
                    if (null != socket && socket.hashCode() == hashCode) {
                        close(CACHE.remove(clientName));
                    }
                }
            }
        }, localDateTimeToDate(LocalDateTime.now().plusMinutes(3L)));

        String targetName = getString(inputStream);
        Socket targetSocket;
        synchronized (CACHE) {
            targetSocket = CACHE.remove(targetName);
            if (null == targetSocket) {
                close(CACHE.put(clientName, clientSocket));
            } else {
                Socket socket = CACHE.remove(clientName);
                if (null != socket && socket.hashCode() != hashCode) {
                    close(socket);
                }
            }
        }

        // 交换地址
        if (null != targetSocket) {
            OutputStream clientOutputStream = clientSocket.getOutputStream();
            OutputStream targetOutputStream = targetSocket.getOutputStream();

            writeString(clientOutputStream, getAddressAndPort(targetSocket));

            writeString(targetOutputStream, getAddressAndPort(clientSocket));
        }
    }

    public static String getString(InputStream inputStream) throws Exception {
        byte[] sLength = getByteArray(inputStream, 4);
        int num = ByteBuffer.wrap(sLength).getInt();
        byte[] sByteArray = getByteArray(inputStream, num);
        return new String(sByteArray);
    }

    public static byte[] getByteArray(InputStream inputStream, int length) throws Exception {
        byte[] result = new byte[length];
        int totalLen = 0;
        do {
            byte[] bytes = new byte[result.length - totalLen];
            int len = inputStream.read(bytes);
            for (int i = 0; i < len; i++) {
                result[totalLen++] = bytes[i];
            }
        } while (totalLen < result.length);
        return result;
    }

    private static void natTypeTest(int serverPort) throws Exception {
        final ServerSocket serverSocket = new ServerSocket(serverPort);
        Thread thread;
        (thread = new Thread(() -> {
            while (!Thread.interrupted()) {
                Socket socket = null;
                try {
                    socket = serverSocket.accept();
                    OutputStream outputStream = socket.getOutputStream();
                    writeString(outputStream, getAddressAndPort(socket));

                } catch (Exception e) {
                    e.printStackTrace();
                } finally {
                    close(socket);
                }
            }
            close(serverSocket);
        })).start();
//        thread.interrupt();
    }

    private static void showAddressAndPort() throws Exception {
        final ServerSocket serverSocket = new ServerSocket(Server.SHOW_ADDRESS_AND_PORT);
        Thread thread;
        (thread = new Thread(() -> {
            while (!Thread.interrupted()) {
                Socket socket = null;
                try {
                    socket = serverSocket.accept();
                    System.out.println(LocalDateTime.now() + " ---> " + getAddressAndPort(socket));
                } catch (Exception e) {
                    e.printStackTrace();
                } finally {
                    close(socket);
                }
            }
            close(serverSocket);
        })).start();
//        thread.interrupt();
    }

    public static void close(Closeable ... closeables) {
        if (null != closeables && closeables.length > 0) {
            for (Closeable closeable : closeables) {
                if (null != closeable) {
                    try {
                        closeable.close();
                    } catch (Exception ignored) { }
                }
            }
        }
    }

    public static Date localDateTimeToDate(LocalDateTime localDateTime) {
        ZoneId zoneId = ZoneOffset.of("+8");
        ZonedDateTime zdt = localDateTime.atZone(zoneId);
        return Date.from(zdt.toInstant());
    }

    public static void writeString(OutputStream outputStream, String s) throws Exception {
        byte[] bytes = s.getBytes();
        outputStream.write(ByteBuffer.allocate(4).putInt(bytes.length).array());
        outputStream.flush();
        outputStream.write(bytes);
        outputStream.flush();
    }

    public static String getAddressAndPort(Socket socket) {
        return socket.getInetAddress().getHostName() + SEPARATOR + socket.getPort();
    }

}
