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
            @RequestParam("mode") Optional<String> mode) {
        Optional<String> chaos = Optional.empty();
        switch (mode.orElse("")) {
            case "1":
                chaos = Optional.of("latency");
                break;
            case "2":
                chaos = Optional.of("cpu");
                break;
            case "3":
                chaos = Optional.of("redis_latency");
                break;
            default:
                break;
        }
        List<User> users = businessService.getUsersCached(chaos, Optional.of(0));
        return ResponseEntity.ok(users);
    }
}