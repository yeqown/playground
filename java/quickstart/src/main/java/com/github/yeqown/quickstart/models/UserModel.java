package com.github.yeqown.quickstart.models;

// import org.bson.types.String;
import org.springframework.data.annotation.Id;

/**
 * UserModel
 */
public class UserModel {

    @Id
    private String id;
    private String username;
    private Integer age;

    public UserModel(String username, Integer age) {
        // this.id = id;
        this.username = username;
        this.age = age;
    }

    /**
     * @return the age
     */
    public Integer getAge() {
        return age;
    }

    /**
     * @param age the age to set
     */
    public void setAge(Integer age) {
        this.age = age;
    }

    /**
     * @return the username
     */
    public String getUsername() {
        return username;
    }

    /**
     * @param username the username to set
     */
    public void setUsername(String username) {
        this.username = username;
    }

    /**
     * @return the id
     */
    public String getId() {
        return id;
    }

    /**
     * @param id the id to set
     */
    public void setId(String id) {
        this.id = id;
    }
}
