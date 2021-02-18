import java.net.InetSocketAddress;
import java.net.Socket;

public class NATTest {

    public static void main(String[] args) throws Exception {
        Socket socket = new Socket(P2P.SERVER_IP, Server.SHOW_ADDRESS_AND_PORT);
    }

    public static void test(int port) throws Exception {
        String s1 = getPublicAddressAndPort(port, Server.NAT_TYPE_TEST_PORT_1);
        Thread.sleep(P2P.SOCKET_BIND_WAIT_TIME);
        String s2 = getPublicAddressAndPort(port, Server.NAT_TYPE_TEST_PORT_2);
        Thread.sleep(P2P.SOCKET_BIND_WAIT_TIME);

        if (s1.equals(s2)) {
            System.out.println("网络地址交换类型验证成功");
        } else {
            throw new RuntimeException("此穿透方法不支持您的网络地址交换类型: " + s1 + " ---> " + s2);
        }

    }

    private static String getPublicAddressAndPort(int localPort, int serverPort) throws Exception {
        Socket socket = new Socket();
        socket.bind(new InetSocketAddress("0.0.0.0", localPort));
        socket.connect(new InetSocketAddress(P2P.SERVER_IP, serverPort));
        String addressAndPort = Server.getString(socket.getInputStream());
        socket.close();
        return addressAndPort;
    }

}
