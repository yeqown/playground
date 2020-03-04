package com.github.yeqown.quickstart.controllers;

import org.springframework.web.bind.annotation.RestController;
import java.util.HashMap;
import java.util.Map;
import javax.servlet.http.HttpServletRequest;
import org.springframework.boot.web.servlet.error.ErrorController;
import org.springframework.web.bind.annotation.ExceptionHandler;
// import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;

@RestController
public class Error implements ErrorController {

    @Override
    public String getErrorPath() {
        return "/error";
    }

    @ExceptionHandler(value=Exception.class)
    @RequestMapping("/error")
    public Map<String, Object> handle(HttpServletRequest request, Exception e) {
        Map<String, Object> map = new HashMap<String, Object>();
        map.put("method", request.getMethod());
        map.put("uri", request.getRequestURI());
        map.put("status", request.getAttribute("javax.servlet.error.status_code"));
        map.put("exception", e.getLocalizedMessage());
        return map;
    }
}
