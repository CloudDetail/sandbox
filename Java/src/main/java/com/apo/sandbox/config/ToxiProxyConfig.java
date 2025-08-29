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
    public ToxiproxyClient toxiproxyClient() {
        String addr = appProperties.getProxyAddr();
        String[] parts = addr.split(":");
        String host = parts[0];
        int port = Integer.parseInt(parts[1]);
        return new ToxiproxyClient(host, port);
    }

    @Bean
    public Proxy toxiProxy(ToxiproxyClient client) {
        try {
            String proxyName = "redis_proxy";
            String redisTarget = appProperties.getRedisHost() + ":" + appProperties.getRedisPort();

            try {
                Proxy existingProxy = client.getProxy(proxyName);
                if (existingProxy != null) {
                    return existingProxy;
                }
            } catch (Exception e) {

            }

            return client.createProxy(proxyName, appProperties.getProxyListenAddr(), redisTarget);
        } catch (Exception e) {
            throw new RuntimeException("Failed to create Toxiproxy proxy for Redis", e);
        }
    }
}