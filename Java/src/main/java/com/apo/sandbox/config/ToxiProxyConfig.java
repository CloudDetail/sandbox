package com.apo.sandbox.config;

import eu.rekawek.toxiproxy.Proxy;
import eu.rekawek.toxiproxy.ToxiproxyClient;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class ToxiProxyConfig {
    @Autowired
    private AppProperties appProperties;

    @Bean
    public Proxy toxiProxy() {
        if (!appProperties.getDeployProxy()) {
            return null;
        }

        try {
            if (!appProperties.getDeployProxy()) {
                return null;
            }
            String addr = appProperties.getProxyAddr();
            String[] parts = addr.split(":");
            String host = parts[0];
            int port = Integer.parseInt(parts[1]);
            ToxiproxyClient client = new ToxiproxyClient(host, port);
            String proxyName = "redis";
            String redisTarget = appProperties.getRedisHost() + ":" + appProperties.getRedisPort();

            Proxy existingProxy = client.getProxy(proxyName);
            if (existingProxy != null) {
                return existingProxy;
            }
            return client.createProxy(proxyName, appProperties.getProxyListenAddr(), redisTarget);
        } catch (Exception e) {
            throw new RuntimeException("Failed to create Toxiproxy proxy for Redis", e);
        }
    }
}