# Plasma - Simple Web Application Firewall (WAF)

Plasma is simple and basic web application firewall.

_Plasma is my personal project and no one should use it in production._

## Features

- [x] Rate-limiting
- [x] SQL Injection protection
- [x] XSS protection
- [x] FAQ endpoint

## FAQ Endpoint

The `/faq` endpoint provides answers to common questions. For example:

```bash
curl "http://localhost:8080/faq?q=movie+about+flat+earth+carried+by+turtle"
```

This will return information about "The Color of Magic" from Terry Pratchett's Discworld series.

## Features I plan to add:

- [ ] Fingerprinting
- [ ] Advanced protection from bots (anomaly detection)
