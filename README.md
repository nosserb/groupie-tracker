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
- Liste de tous les artistes disponibles (affiché en grille sur 2 colonnes)
- **Système de filtres avancé et optimisé** avec 4 critères principaux:

#### Filtres de plage (Range Filters)
  - **Année de création** : Utilise des sliders interactifs pour sélectionner une plage min-max
    - Entrées numériques directes ou sliders visuels
    - Affichage en temps réel de la plage sélectionnée
    - Plage disponible : 1950-2030
  - **Date du premier album** : Même système de sliders que la création
    - Sélection précise de la période album souhaitée
    - Plage disponible : 1950-2030

#### Filtres par cases à cocher (Checkbox Filters)
  - **Nombre de membres** : Filtrage multi-sélection
    - Options : 1, 2, 3, 4, 5, 6 membres et 7+ membres
    - Permet de combiner plusieurs sélections
  - **Lieux de concerts** : Liste dynamique de tous les lieux disponibles
    - Généré automatiquement à partir des données
    - Tri alphabétique
    - Filtrage par localisation avec support des variantes (Seattle, Washington USA)

#### Recherche avancée
  - Barre de recherche en temps réel
  - Recherche dans les noms d'artistes, membres et lieux
  - Combinaison possible avec tous les filtres
  - Historique de recherche via URL (paramètre `?search=`)

#### Gestion des filtres
  - Synchronisation automatique des sliders et champs numériques
  - Application instantanée des filtres (sans rechargement de page)
  - Bouton "Réinitialiser les filtres" pour nettoyer tous les critères
  - Performance optimisée avec API backend dédié

- Clic sur un artiste pour voir les détails complets

### API Filters
- `GET /api/filter` - Endpoint d'API pour le filtrage avancé
  - Paramètres : `creationDateMin`, `creationDateMax`, `firstAlbumMin`, `firstAlbumMax`, `memberCounts`, `locations`
  - Retourne la liste filtrée des artistes en JSON
  - Performance : Filtrage côté serveur pour optimiser les performances

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
- `GET /api/filter` - Filtrage avancé des artistes
  - Paramètres query string :
    - `creationDateMin` : Année minimale de création (entier)
    - `creationDateMax` : Année maximale de création (entier)
    - `firstAlbumMin` : Année minimale du premier album (entier)
    - `firstAlbumMax` : Année maximale du premier album (entier)
    - `memberCounts` : Nombres de membres à filtrer (liste d'entiers séparés par virgule)
    - `locations` : Lieux à filtrer (liste de localisations séparées par virgule)
  - Retourne un objet JSON avec structure:
    ```json
    {
      "artists": [...],
      "total": 15
    }
    ```
- `GET /artist/{id}` - Récupère les détails complets d'un artiste (id: 1-52)
  - Retourne l'artiste, ses lieux de concerts, dates et relations
- `GET /api/search` - Recherche multi-critères (artistes, membres, lieux, dates)
  - Paramètre : `q` (query string)
  - Retourne une liste de résultats avec suggestion de recherche
- `GET /api/geocode` - Service de géocodage pour localiser les lieux
  - Paramètre : `address` (adresse à géocoder)
  - Retourne les coordonnées latitude/longitude

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
- **Système de filtrage optimisé** :
  - Endpoint API `/api/filter` pour filtrage multi-critères côté serveur
  - Support des filtres de plage (création date, album date)
  - Support des filtres multi-sélection (membres, lieux)
  - Combinaison de plusieurs filtres simultanément
  - Extraction intelligente des années de dates (format "DD-MM-YYYY")
  - Filtrage de localisation sensible aux variantes de noms
- Support des modales dynamiques via JavaScript
- Cache de géocodage pour optimiser les appels API

### Frontend
- Requêtes AJAX/Fetch pour récupérer les données d'artistes
- Affichage dynamique des informations via JavaScript vanilla
- **Système de filtrage avancé** :
  - Sliders interactifs avec synchronisation double-sens (slider ↔ input numérique)
  - Affichage en temps réel de la plage sélectionnée
  - Filtres à cases à cocher multi-sélection
  - Application instantanée des filtres sans rechargement de page
  - Optimisation API : requête unique au serveur avec tous les paramètres
- Système de recherche en temps réel côté client
- Modales interactives pour les détails d'artistes
- Design responsive pour mobile et desktop
- Gestion des erreurs avec messages informatifs
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

type FilterResponse struct {
    Artists []Artist // Artistes filtrés
    Total   int      // Nombre total d'artistes retournés
}
```

### Tests unitaires
- Tests du système de filtrage :
  - `TestExtractYear` : Extraction correcte des années depuis les dates
  - `TestFilterArtists` : Filtrage multi-critères (création date, album date, membres, lieux)
  - `TestFilterHandler` : API endpoint de filtrage
- Tests des endpoints existants (fetch JSON, handlers, geocoding)
- Couverture : 100% des cas de test critiques

Pour exécuter les tests :
```bash
go test -v
```

### Gestion des erreurs
- Le serveur gère les erreurs HTTP (404, 500, etc.)
- Logging des erreurs dans la console du serveur
- Gestion gracieuse des connexions API défaillantes
- Affichage des messages d'erreur côté client
- Validation des paramètres de requête
- Gestion des edge cases (données manquantes, formats invalides)

### Bonnes pratiques
- Code Go organisé et bien structuré
- Séparation claire des responsabilités (frontend/backend)
- Utilisation de templates Go pour le rendu HTML sécurisé
- Gestion appropriée des ressources (defer, fermeture des connexions)
- Absence de dépendances externes - utilisation uniquement des packages standard Go
- Code JavaScript vanilla sans frameworks lourds
- Filtrage côté serveur pour sécurité et performance
- Synchronisation bidirectionnelle pour les sliders (UX améliorée)

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
