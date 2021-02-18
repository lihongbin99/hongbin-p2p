public class Main {

    public static void main(String[] args) {
        new Thread(Client1::main).start();
        new Thread(Client2::main).start();
    }

}
