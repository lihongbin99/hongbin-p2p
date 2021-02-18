import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;

public class P2PTest {

    public static void test(final Socket socket, String targetName) throws Exception {
        // TODO 成功连接与对方的 Socket 连接后即可传输数据
        new Thread(() -> {
            try {
                InputStream inputStream = socket.getInputStream();
                byte[] bytes = new byte[1024 * 64];
                int len;
                while ((len = inputStream.read(bytes)) != -1) {
                    System.out.println(new String(bytes, 0, len));
                }
            } catch (Exception e) {
                e.printStackTrace();
            }
        }).start();

        OutputStream outputStream = socket.getOutputStream();
        outputStream.write(("如果在" + targetName + "看到此信息则说明成功").getBytes());
        outputStream.flush();

        System.in.read();
    }

}
