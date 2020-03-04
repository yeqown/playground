package com.github.yeqown.quickstart.controllers;

/**
 * CommonResponse
 */
public class CommonResponse {

    public Integer code;
    public String message;

    CommonResponse() {
        this.code = 0;
        this.message = "OK";
    }

    CommonResponse(Integer code, String msg) {
        this.code = code;
        this.message = msg;
    }

    /**
     * @param code the code to set
     */
    public void setCode(Integer code) {
      this.code = code;
    }

    /**
     * @param message the message to set
     */
    public void setMessage(String message) {
      this.message = message;
    }
}