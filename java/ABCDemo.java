abstract class ABC {
    private int Code = 0;

    public void Func1() {
        this.Func2();
    }

    private void Func2() {
        System.out.println("Func2 called: " + this.Code);
    }
}

class NewABC extends ABC {
    private int Code = 2;

    public NewABC() {
        System.out.println("called from constructor start");
        this.Func2();
        System.out.println("called from constructor end");
    }

    /**
     * @return the code
     */
    public int getCode() {
        return Code;
    }

    /**
     * @param code the code to set
     */
    public void setCode(int code) {
        this.Code = code;
    }

    private void Func2() {
        super.Func1();
    }
}

/**
 * ABCDemo
 */
public class ABCDemo {
    public static void main(String[] args) {
        NewABC nabc = new NewABC();
        nabc.Func1();
    }
}

/** Output:
 * 
 * called from constructor start
 * Func2 called: 0
 * called from constructor end
 * Func2 called: 0
 * 
 */