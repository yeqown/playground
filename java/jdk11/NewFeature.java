import java.math.BigDecimal;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.function.Consumer;

import javax.script.ScriptEngine;
import javax.script.ScriptEngineManager;

/**
 * NewFeature
 */
public class NewFeature {
    public static void main(String[] args) {
        // strings
        var chars = " CHINA!!! ";
        System.out.println("chars.repaeat(3) = " + chars.repeat(3));
        System.out.println("chars.isBlank() = " + chars.isBlank());
        System.out.println("chars.strip() = " + chars.strip());

        // lambdas
        Consumer<Integer> moneyConsumer = (Integer money) -> {
            System.out.println("i got this much money = " + money);
        };
        moneyConsumer.accept(10000);

        // http client
        try {
            HttpClient client = HttpClient.newBuilder().build();
            HttpRequest request = HttpRequest.newBuilder().GET().uri(URI.create("https://www.baidu.com")).build();
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            System.out.println(response.body());
        } catch (Exception e) {
            e.printStackTrace();
        }

        // unicode 10
        System.out.println("\u20BF");

        // javascriptEngine nashorn
        try {
            ScriptEngine engine = new ScriptEngineManager().getEngineByName("nashorn");
            engine.eval("print('help, i am looking')");
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}