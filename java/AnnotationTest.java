import java.lang.annotation.Documented;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Annotation;

@Retention(RetentionPolicy.RUNTIME)
@Documented // this is meta annotation to annotate the MyAnnotation
@interface MyAnnotation {
    public String msg() default "no_msg";
}

@MyAnnotation(msg = "hahah")
class Demo {
    public Demo() {
        MyAnnotation annotation = this.getClass().getAnnotation(MyAnnotation.class);
        System.out.println("Demo constructor: " + annotation.toString() + " has attr mag=" + annotation.msg());
    }
}

/**
 * AnnotationTest
 */

public class AnnotationTest {

    public static void main(String[] args) {
        Demo obj = new Demo();
        // Class classz = obj.getClass();
        Annotation[] annotations = obj.getClass().getAnnotations();

        for (Annotation a : annotations) {
            System.out.println(a.toString());
        }
    }
}