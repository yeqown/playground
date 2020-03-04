package com.github.yeqown.quickstart.repository;

import com.github.yeqown.quickstart.models.UserModel;
import org.bson.types.ObjectId;
import org.springframework.data.mongodb.repository.MongoRepository;



/**
 * UserRepository
 */
public interface UserRepository extends MongoRepository<UserModel, ObjectId> {

    UserModel findUserById(ObjectId Id);
    UserModel findUserByUsername(String username);
}