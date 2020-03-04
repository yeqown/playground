package com.github.yeqown.quickstart.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * Index
 */
@RestController
public class Index {

    @GetMapping("/")
    public String handle() {
        return "this is index entry";
    }    
}