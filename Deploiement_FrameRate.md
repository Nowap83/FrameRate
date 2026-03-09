# Déploiement et Architecture de l'Application FrameRate

Ce document récapitule l'architecture de déploiement de l'application FrameRate et les procédures mises en place pour assurer un déploiement continu, fiable et sécurisé.

## 1. Description de l'Architecture de Déploiement

L'application FrameRate est composée de plusieurs services conteneurisés via Docker, orchestrés par Docker Compose. Le déploiement s'effectue sur un serveur Proxmox hébergeant un conteneur LXC (Linux Containers).

L'architecture comprend deux environnements sur le même serveur LXC :
- **Environnement de Pré-production (`preprod`)** : Déployé à partir de la branche `preprod` pour valider les nouvelles fonctionnalités avant leur passage en production.
- **Environnement de Production (`main`)** : Déployé à partir de la branche `main`, constituant la version stable accessible aux utilisateurs finaux.

### Composants de l'Application :
- **Frontend** : Application React + Vite, servie par un serveur web Nginx.
- **Backend** : API REST développée en Go.
- **Base de Données** : PostgreSQL (versions 16).
- **Cache** : Redis (version 7).

### Gestion de l'Accès Distant :
Afin de ne pas exposer directement le serveur LXC sur Internet, nous utilisons un **Tunnel Cloudflare (cloudflared)**. Cela permet :
1. D'acheminer le trafic HTTP/HTTPS entrant de manière sécurisée vers les ports locaux appropriés.
2. D'établir une connexion SSH sécurisée pour les GitHub Actions sans avoir à ouvrir le port 22 sur un pare-feu public.

## 2. Évolutions Complexes (Ports et Conflits)

Puisque les deux environnements (`preprod` et `main`) cohabitent sur le même serveur LXC, le déploiement du second environnement (la production) nécessite des adaptations pour éviter les conflits de ports et de noms de conteneurs.

Pour la **Pré-production** (`docker-compose.yml`), les ports par défaut sont utilisés :
- Frontend : `80`
- Backend : `8080`
- PostgreSQL : `5432`
- Redis : `6379`

Pour la **Production** (`docker-compose.prod.yml`), nous avons mis à jour la configuration avec des ports distincts :
- Frontend : `81` (C'est ce port local qui est exposé via le tunnel Cloudflare pour le domaine de production)
- Backend : `8081`
- PostgreSQL : `5433`
- Redis : `6380`

Des préfixes/suffixes adaptés ont également été ajoutés aux noms des conteneurs (ex: `framerate_backend_prod`) et des volumes.

## 3. Enjeux de Sécurité

La procédure de déploiement prend en compte plusieurs aspects critiques de la sécurité :

- **Gestion des Secrets et Variables d'Environnement** :
  - Les clés sensibles (mots de passe BDD, clés JWT, clés API externes) et les identifiants de connexion SSH (`SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY`) ne sont jamais versionnés dans le code.
  - Ils sont stockés de manière chiffrée dans **GitHub Secrets**.
  - Ces secrets sont injectés dynamiquement sous forme de variables d'environnement (`.env`) au moment du déploiement ou passés au conteneur Docker.
- **Accès Distants Sécurisés** :
  - L'accès SSH utilisé par GitHub Actions pour se connecter au serveur LXC ne se fait pas via le réseau public, mais passe par l'outil CLI `cloudflared` configuré en tant que `ProxyCommand`. L'IP réelle du serveur n'est donc pas exposée.
- **Chiffrement et Certificats** :
  - L'accès public à l'application est protégé par Cloudflare, qui fournit automatiquement des certificats SSL/TLS valides, garantissant le chiffrement des données de bout en bout (HTTPS) entre l'utilisateur et Cloudflare, puis entre Cloudflare et le tunnel sur le LXC.

## 4. Procédure de Déploiement Continues (CI/CD)

Le déploiement est entièrement automatisé via des workflows **GitHub Actions**.

### 1. Procédure de vérification des dépendances et de l'état des versions :
- **Etape de "Checkout"** : Le code source est récupéré de la branche ciblée.
- **Environnements d'exécution** : Le workflow configure des environnements stricts (`setup-go@v5` pour Go 1.24 et `setup-node@v4` pour Node.js 20). Ainsi, le déploiement utilise systématiquement les mêmes versions d'outils que celles définies dans l'architecture, évitant l'effet "ça marche sur ma machine".
- **Installation des dépendances** : Lors de la construction des images Docker, les dépendances Frontend (`npm ci`) et Backend (`go mod download`) sont fraîchement téléchargées pour s'assurer de leur intégrité.

### 2. Procédure d'exécution des tests (Intégration et Acceptation)
Le job `test` est un pré-requis strict au déploiement (`needs: test`).
- **Tests Backend (Système et Intégration)** : Les tests unitaires et les tests des services Go sont exécutés (`go test -v ./...`).
- **Tests Frontend** : La suite de tests React/Vite est lancée (`npm run test:run`).
Si l’un des environnements de test échoue, le processus s'arrête et **le déploiement en production ou en pré-production est annulé**. L'équipe de développement en est alors immédiatement informée via les notifications de l'interface GitHub.

### 3. Procédure de Déploiement :
Si les tests sont validés, le job `deploy` est exécuté :
1. Installation de `cloudflared` sur le runner GitHub.
2. Configuration de l'agent SSH avec la clé privée sécurisée provenant des secrets GitHub.
3. Exécution d'un script bash dynamique via SSH sur le serveur cible. Ce script va :
   - Créer le dossier cible (`/opt/main` ou `/opt/preprod`).
   - Initialiser ou mettre à jour un dépôt git local branché sur la bonne révision.
   - Forcer la synchronisation avec le code validé (`git reset --hard origin/<nom_branche>`).
   - Redémarrer gracieusement les services applicatifs via Docker Compose (`docker compose down && docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build`).

