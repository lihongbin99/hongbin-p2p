import java.net.Socket;

public class Client1 {

    private final static int PORT = 12580;
    private final static String CLIENT_NAME = "client1";
    private final static String TARGET_NAME = "client2";

    public static void main(String[] args) {
        main();
    }

    public static void main() {
        try {
            Socket socket = P2P.getSocket(Client1.CLIENT_NAME, Client1.TARGET_NAME, Client1.PORT);
            P2PTest.test(socket, Client1.TARGET_NAME);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
