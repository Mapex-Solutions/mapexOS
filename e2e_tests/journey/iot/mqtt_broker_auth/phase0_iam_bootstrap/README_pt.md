# Fase 0 — Bootstrap de IAM (login do admin seed)

## O que este teste prova

O usuário admin seed consegue fazer login via mapexIam e o JWT
resultante carrega cobertura de org-context para a organização raiz
seed. O ClientSet propaga esse bearer + header `X-Org-Context` para
todos os services, de modo que qualquer fase subsequente possa operar
como esse admin sem refazer login.

A fase:

1. Faz POST com credenciais do admin (`admin@mapex.local`) em `/auth/login`; publica JWT + organizationID no bag.
2. Verifica que o JWT é estruturalmente válido (parse, assinatura, expiração no futuro).
3. Verifica que o JWT carrega acesso a `MapexosOrgID` via endpoint de coverage do IAM.

## Como rodar

```bash
cd e2e_tests
go test -tags=saga -count=1 ./journey/iot/mqtt_broker_auth/phase0_iam_bootstrap/
```

## Requisitos

- Serviço `mapexos` ativo em `MAPEXOS_URL` (default `http://localhost:5000`).
- Usuário admin seed (`admin@mapex.local` / `mapex@123`) provisionado pelo seed canônico do mongodb-init.
- Job de build de coverage já produziu a entrada wildcard `mapex.*` para o admin seed ancorada na org raiz seed.
