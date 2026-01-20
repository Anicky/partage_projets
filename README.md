# Partage de projets

Une plateforme de partage de projets, développée en Go avec le framework Gin.
Ce projet permet aux utilisateurs de partager leurs projets, de les consulter, et d'échanger via des commentaires.

## Fonctionnalités

- **Gestion des utilisateurs**
  - Inscription d'un utilisateur
  - Connexion d'un utilisateur
- **Gestion des projets**
  - Création d'un projet
  - Modification d'un projet
  - Suppression d'un projet
  - Affichage de tous les projets
  - Affichage d'un projet
  - Ajout / suppression d'un like sur un projet
- **Commentaires**
  - Ajout d'un commentaire sur un projet

## Installation et configuration

### Prérequis

- **Go** (>= 1.25.0)
- **PostgresL**

### Configuration des variables d'environnement

Créer un fichier `.env` à la racine du projet, en reprenant le contenu du fichier `.env.dist`, et en le personnalisant avec vos informations.

### Lancement de l'application

```bash
go run main.go
```

Le serveur démarrera par défaut sur `http://localhost:8080`.

## Documentation

### Swagger

Une fois le serveur lancé, vous pouvez accéder à la documentation Swagger à l'adresse suivante :
`http://localhost:8080/swagger/index.html`

### Postman

Vous pouvez importer la collection Postman en utilisant le fichier `postman_collection.json` placé à la racine du projet.