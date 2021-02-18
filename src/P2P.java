import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;

public class P2P {

//    public final static String SERVER_IP = "0.0.0.0";
    public final static String SERVER_IP = "149.129.57.253";

    public final static int SOCKET_BIND_WAIT_TIME = 3/*s*/ * 1000/*ms*/;

    public static Socket getSocket(String clientName, String targetName, int localPort) throws Exception {
        // 已知网络交换地址为 非Symmetric 或曾经使用此代码验证过一次后即可注销此代码
        NATTest.test(localPort);

        Socket socket = new Socket();
        socket.bind(new InetSocketAddress("0.0.0.0", localPort));
        socket.connect(new InetSocketAddress(SERVER_IP, Server.SERVER_PORT));

        OutputStream outputStream = socket.getOutputStream();
        Server.writeString(outputStream, clientName);
        Server.writeString(outputStream, targetName);

        String addressAndPort = Server.getString(socket.getInputStream());
        socket.close();

        String targetIp;
        int targetPort;
        try {
            String[] split = addressAndPort.split(Server.SEPARATOR);
            targetIp = split[0];
            targetPort = Integer.parseInt(split[1]);
        } catch (Exception e) {
            throw new RuntimeException(addressAndPort);
        }

        System.out.println("target_ip    = " + targetIp);
        System.out.println("target_port  = " + targetPort);
        Thread.sleep(SOCKET_BIND_WAIT_TIME);

        socket = new Socket();
        socket.bind(new InetSocketAddress("0.0.0.0", localPort));
        System.out.println("开始连接 ---> " + targetIp + ":" + targetPort);
        socket.connect(new InetSocketAddress(targetIp, targetPort), 60000);
        return socket;

    }

}
