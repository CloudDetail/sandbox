package com.apo.sandbox.controller;

import com.apo.sandbox.config.AppProperties;
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
    private final AppProperties appProperties;

    public BusinessController(BusinessService businessService, AppProperties appProperties) {
        this.businessService = businessService;
        this.appProperties = appProperties;
    }

    @GetMapping("/users/1")
    public ResponseEntity<List<User>> getUsersWithLatency(@RequestParam("mode") Optional<String> mode) {
        int duration = appProperties.getLatencyFaultDefaultDelay();
        List<User> users = businessService.getUsersWithLatency(mode, duration);
        return ResponseEntity.ok(users);
    }

    @GetMapping("/users/2")
    public ResponseEntity<List<User>> getUsersWithCPUBurn(@RequestParam("mode") Optional<String> mode) {
        int duration = appProperties.getCpuFaultDefaultDuration();
        List<User> users = businessService.getUsersWithCPUBurn(mode, duration);
        return ResponseEntity.ok(users);
    }

    @GetMapping("/users/3")
    public ResponseEntity<List<User>> getUsersWithRedisLatency(@RequestParam("mode") Optional<String> mode) {
        int duration = appProperties.getRedisFaultDefaultDelay();
        List<User> users = businessService.getUsersWithRedisLatency(mode, duration);
        return ResponseEntity.ok(users);
    }
}