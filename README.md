# Groupie Tracker

Un projet web qui récupère et affiche les informations de groupes musicaux et artistes via l'API Groupie Trackers.

## Description

Groupie Tracker est une application web développée en Go qui permet d'afficher des informations sur divers artistes musicaux, incluant :
- Nom du groupe/artiste
- Image du groupe
- Année de création
- Date du premier album
- Membres du groupe
- Lieux de concerts
- Dates de concerts
- Relations entre dates et lieux

## Technologies utilisées

- **Backend** : Go 1.25.3 (packages standard uniquement)
- **Frontend** : HTML5, CSS3, JavaScript vanilla
- **API externe** : https://groupietrackers.herokuapp.com/api
- **Serveur HTTP** : Go http package
- **Architecture** : MVC (Model-View-Controller) pattern

## Structure de l'API

L'API est composée de 4 parties :

1. **Artists** : Informations sur les groupes et artistes
2. **Locations** : Lieux de concerts (passés et à venir)
3. **Dates** : Dates de concerts (passés et à venir)
4. **Relation** : Lie les données des artistes, dates et lieux

## Installation et utilisation

### Prérequis

- Go 1.25.3 ou supérieur

### Lancement du serveur

1. Clonez le dépôt :
```bash
git clone <repository-url>
cd groupie-tracker
```

2. Lancez le serveur :
```bash
go run main.go
```

3. Ouvrez votre navigateur à l'adresse : `http://localhost:8080`

**Note** : Le serveur se connecte automatiquement à l'API Groupie Trackers au démarrage pour récupérer les données.

## Fonctionnalités

### Page d'accueil (/)
- Affiche un podium avec les 3 premiers artistes
- Navigation vers la page complète des artistes
- Design responsif et moderne

### Page Artists (/artists.html)
- Liste de tous les artistes disponibles (affiché en grille)
- **Système de filtres avancé** avec 4 critères:
  - **Filtre par nombre de membres** (checkbox multi-sélection: 1, 2, 3, 4, 5, 6, 7+)
  - **Filtre par année de création** (range checkbox: 1980-1990, 1991-2000, 2001-2010, 2011-2020, 2021-2030)
  - **Filtre par date du premier album** (range checkbox: 1960-1980, 1981-1995, 1996-2005, 2006-2015, 2016-2025)
  - **Filtre par lieux de concerts** (checkbox multi-sélection dynamique de tous les lieux)
- Recherche en temps réel d'artistes
- Combinaison de plusieurs filtres possibles
- Bouton "Réinitialiser les filtres" pour nettoyer tous les filtres
- Clic sur un artiste pour voir les détails complets

### Modal d'information
- Affiche les détails complets d'un artiste :
  - Image de haute qualité
  - Année de création
  - Date du premier album
  - Liste complète des membres
  - Lieux de concerts (avec possibilité de les afficher sur une carte)
  - Dates de concerts
  - Relations entre dates et lieux
- Fermeture en cliquant sur le bouton X ou en dehors de la modal

### Page Crédits (/credits.html)
- Informations sur les contributeurs et sources

### Page Carte (/map.html)
- Affichage interactif des lieux de concerts sur une carte
- Intégration avec les données de localisation des artistes

## Structure du projet

```
groupie-tracker/
├── main.go                      # Serveur principal et logique backend
├── main_test.go                 # Tests unitaires
├── index.html                   # Page d'accueil
├── artists.html                 # Page liste des artistes
├── credits.html                 # Page crédits
├── map.html                     # Page carte interactive
├── index.css                    # Styles page d'accueil
├── artists.css                  # Styles page artistes
├── credits.css                  # Styles page crédits
├── modal.css                    # Styles pour les modals d'artistes
├── map.css                      # Styles pour la carte
├── go.mod                       # Module Go (v1.25.3)
├── package.json                 # Configuration npm (optionnel)
├── server.js                    # Alternative serveur Node.js (optionnel)
├── README.md                    # Documentation du projet
└── public/                      # Ressources statiques
    ├── images/                  # Images du site
    └── font/                    # Polices de caractères
```

## API Endpoints

### Pages HTML
- `GET /` - Page d'accueil (podium des top 3)
- `GET /artists.html` - Page liste complète des artistes
- `GET /credits.html` - Page crédits
- `GET /map.html` - Page carte interactive

### API JSON
- `GET /api/artists` - Récupère tous les artistes en JSON
- `GET /artist/{id}` - Récupère les détails complets d'un artiste (id: 1-52)
  - Retourne l'artiste, ses lieux de concerts, dates et relations

### Fichiers statiques
- `GET /public/{chemin}` - Ressources statiques (images, fonts, etc.)
- `GET /*.css` - Fichiers de styles CSS

## Fonctionnalités techniques

### Backend (Go)
- Récupération de données depuis l'API Groupie Trackers au démarrage
- Gestion des routes HTTP avec le package standard `net/http`
- Templating HTML dynamique avec le package `html/template`
- Marshaling/Unmarshaling JSON avec le package `encoding/json`
- Structures de données typées pour les artistes, lieux et dates
- Gestion des erreurs robuste avec logging
- Support des modales dynamiques via JavaScript

### Frontend
- Requêtes AJAX/Fetch pour récupérer les données d'artistes
- Affichage dynamique des informations via JavaScript vanilla
- Système de recherche en temps réel côté client
- Modales interactives pour les détails d'artistes
- Design responsive pour mobile et desktop
- Intégration de carte (possibilité d'utiliser Leaflet, Google Maps, etc.)

### Structures de données
```go
type Artist struct {
    ID           int      // Identifiant unique
    Image        string   // URL de l'image
    Name         string   // Nom du groupe/artiste
    Members      []string // Membres du groupe
    CreationDate int      // Année de création
    FirstAlbum   string   // Date du premier album
    Locations    string   // URL de l'endpoint locations
    ConcertDates string   // URL de l'endpoint dates
    Relations    string   // URL de l'endpoint relations
}
```

### Gestion des erreurs
- Le serveur gère les erreurs HTTP (404, 500, etc.)
- Logging des erreurs dans la console du serveur
- Gestion gracieuse des connexions API défaillantes
- Affichage des messages d'erreur côté client

### Bonnes pratiques
- Code Go organisé et bien structuré
- Séparation claire des responsabilités (frontend/backend)
- Utilisation de templates Go pour le rendu HTML sécurisé
- Gestion appropriée des ressources (defer, fermeture des connexions)
- Absence de dépendances externes - utilisation uniquement des packages standard Go
- Code JavaScript vanilla sans frameworks lourds

## Points d'amélioration possibles

- Implémentation d'un système de cache pour les données API
- Ajout de tests d'intégration supplémentaires
- Optimisation des performances de recherche
- Implémentation d'une base de données locale (SQLite)
- Responsive design amélioré
- Animations CSS pour une meilleure UX
- Déploiement sur un serveur (Heroku, AWS, etc.)

## Auteurs

- gcouvri
- vbosson
- kcamesel

## Licence

Ce projet est développé dans le cadre d'un exercice pédagogique à Zona.

## Ressources utiles

- [Documentation Go](https://golang.org/doc/)
- [API Groupie Trackers](https://groupietrackers.herokuapp.com/api)
- [Documentation MDN Web](https://developer.mozilla.org/fr/)
- [Guide Markdown](https://www.markdownguide.org/)
