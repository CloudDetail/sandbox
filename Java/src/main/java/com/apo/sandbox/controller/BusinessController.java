package com.apo.sandbox.controller;

import com.apo.sandbox.model.User;
import com.apo.sandbox.service.BusinessService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/api")
public class BusinessController {

    private final BusinessService businessService;

    public BusinessController(BusinessService businessService) {
        this.businessService = businessService;
    }

    @GetMapping("/users")
    public ResponseEntity<List<User>> getUsersCached(
            @RequestParam("chaos") Optional<String> chaos,
            @RequestParam("duration") Optional<Integer> duration) {

        List<User> users = businessService.getUsersCached(chaos, duration);
        return ResponseEntity.ok(users);
    }
}