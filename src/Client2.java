import java.net.Socket;

public class Client2 {

    private final static int PORT = 22580;
    private final static String CLIENT_NAME = "client2";
    private final static String TARGET_NAME = "client1";

    public static void main(String[] args) {
        main();
    }

    public static void main() {
        try {
            Socket socket = P2P.getSocket(Client2.CLIENT_NAME, Client2.TARGET_NAME, Client2.PORT);
            P2PTest.test(socket, Client2.TARGET_NAME);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

}
