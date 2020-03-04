/**
 * InnerInterface
 */
interface InnerInterface {
    void IFunc(String s);
}

class MyInterface implements InnerInterface {
    private String d1;

    public MyInterface(String s) {
        this.d1 = s;
    }

    @Override
    public void IFunc(String s) {
        System.out.println(this.d1 + " " + s);
    }
}

class MyInterface2 implements InnerInterface {

    @Override
    public void IFunc(String s) {
        System.out.println(String.format("%s is recieved", s));
    }
}

/**
 * Interface
 */
public class Interface {

    public static void main(String[] args) {
        InnerInterface i1 = new MyInterface("origin");
        InnerInterface i2 = new MyInterface2();

        i1.IFunc("second");
        i2.IFunc("second");
    }
}