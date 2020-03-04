package com.github.yeqown.quickstart.controllers;

import java.util.ArrayList;
import java.util.List;
import javax.validation.Valid;
import javax.validation.constraints.Max;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import javax.validation.constraints.Min;
import com.github.yeqown.quickstart.models.UserModel;
import com.github.yeqown.quickstart.repository.UserRepository;
// import org.apache.logging.log4j.message.Message;
// import org.bson.conversions.Bson;
import org.bson.types.ObjectId;
// import org.hibernate.validator.constraints.Length;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
// import org.springframework.data.domain.Page;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ModelAttribute;
// import org.springframework.web.bind.annotation.ModelAttribute;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

class UserListResponse extends CommonResponse {
    public List<UserModel> users;
    public Long total;

    UserListResponse() {
        this.users = new ArrayList<UserModel>();
        // this.total = 0;
    }

    UserListResponse(Integer code, String message) {
        super(code, message);
        this.users = new ArrayList<UserModel>();
        // this.total = 0;
    }
}


class userRegisterForm {
    @NotNull
    @Size(min = 1, max = 20)
    private String name;

    @NotNull
    @Min(1)
    @Max(120)
    private Integer age;

    /**
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
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
}


class Filter {
    public Integer limit = 10;
    public Integer skip = 0;
}


@RestController
public class UserController {

    @Autowired
    private UserRepository repository;

    private static final Logger logger = LoggerFactory.getLogger(UserController.class);

    @GetMapping("/users")
    public UserListResponse getUsers(@ModelAttribute Filter filter) {
        var resp = new UserListResponse();
        resp.total = this.repository.count();
        resp.users = this.repository.findAll();
        return resp;
    }

    @PostMapping("/user")
    public CommonResponse registerUser(@Valid userRegisterForm form, BindingResult bindingResult) {
        var resp = new CommonResponse();
        if (bindingResult.hasErrors()) {
            logger.info("errors: %v", bindingResult.getAllErrors());
            resp.setCode(-1);
            resp.setMessage(bindingResult.toString());
            return resp;
        }
        logger.info("request form: %v", form);
        this.repository.insert(new UserModel(form.getName(), form.getAge()));
        return resp;
    }

    @DeleteMapping("/user/{id}")
    public CommonResponse deleteUser(@PathVariable String id) {
        var resp = new CommonResponse();
        repository.deleteById(new ObjectId(id));
        return resp;
    }
}
