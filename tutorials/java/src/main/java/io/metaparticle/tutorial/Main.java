package io.metaparticle.tutorial;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import com.sun.net.httpserver.HttpServer;

public class Main {
    private static final int port = 8080;

    public static void main(String[] args) {
    	try {
       		HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
                server.createContext("/", new HttpHandler() {
                    @Override
                    public void handle(HttpExchange t) throws IOException {
                        String msg = "Hello Containers [" + t.getRequestURI() + "] from " + System.getenv("HOSTNAME");
                        t.sendResponseHeaders(200, msg.length());
                        OutputStream os = t.getResponseBody();
                        os.write(msg.getBytes());
                        os.close();
                        System.out.println("[" + t.getRequestURI() + "]");
                    }
                });
                server.start();
        } catch (IOException ex) {
                ex.printStackTrace();
        }
    }
}
