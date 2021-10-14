# [AMIGA-7] Web de sellado de tiempo

Este repositorio configura un servidor de sellado de tiempo con un frontend web
para el proyecto AMIGA-7. El software [utsa](https://github.com/varrrro/utsa) es
utilizado como servidor de sellado según el estándar RFC 3161.

## Requisitos

- Docker
- Ansible
- OpenSSL

## Descarga y uso

Para usar este sistema de sellado, debemos clonar el repositorio teniendo en
cuenta que contiene
[submódulos](https://git-scm.com/book/en/v2/Git-Tools-Submodules), de forma que
usaremos el argumento `--recurse-submodules` al realizar el `git clone`.

Antes de poder levantar el sistema, debemos tener listos los certificados a
usar para el proceso de sellado. Para facilitar la prueba del sistema, se
incluye el script `./scripts/create_certs.sh` que genera unos certificados en
la carpeta `certs`.

Se incluye también un _playbook_ de Ansible que despliega el sistema al
completo con solo ejecutar el siguiente comando:

```
ansible-playbook deploy/website.yml
```

Una vez hecho esto, la web de sellado debería estar accesible en `http://localhost:8000`.
