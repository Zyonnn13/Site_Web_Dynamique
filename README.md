# Site Web Dynamique (Go)

Application web simple en Go qui affiche une liste de produits, une page de détails, et permet d'ajouter un produit via un formulaire (avec upload d'image optionnel et fallback automatique si l'image est invalide).

## Prérequis

- Go 1.20+ installé
- Windows PowerShell (les chemins ci-dessous sont pour Windows)

## Structure du projet

```
assets/
	css/
		style.css
	img/
		logo/
		products/
			00958058-veste-lazare-gris-moyen_5_1.avif   # image par défaut/fallback
	temp/
		add.html      # formulaire d'ajout
		index.html    # liste des produits
		product.html  # détails produit
main.go
go.mod
README.md
```

## Lancer le serveur

Dans un terminal PowerShell à la racine du projet:

```powershell
go run .
```

Par défaut, le serveur écoute sur `http://localhost:8080`.

## Routes

- `GET /` — accueil avec la liste des produits
- `GET /product/{id}` — page de détails d’un produit
- `GET /add` — formulaire d’ajout d’un produit
- `POST /add/create` — traitement du formulaire (validation + ajout en mémoire + redirection)
- `GET /static/...` — fichiers statiques (CSS, images)

## Templates et rendu

- `index.html` parcourt la liste globale `products` et affiche chaque carte produit. Si un produit a un prix réduit (`ReducPrice`), l'ancien prix est barré en blanc et le nouveau prix est en rouge.
- `product.html` affiche les détails complets d’un produit sélectionné.
- `add.html` est le formulaire d’ajout (nom, prix, prix réduit optionnel, description, stock, tailles, et image optionnelle).

## Ajouter un produit

1. Ouvre `http://localhost:8080/add`.
2. Renseigne les champs obligatoires: nom, prix, description, stock.
3. Facultatif: prix réduit, tailles, et une image (png/jpg/jpeg/webp/avif/gif).
4. Envoie le formulaire: le produit est ajouté en mémoire et tu es redirigé vers `/product/{id}`.

### Gestion des images

- Si tu choisis une image, elle est enregistrée sous `assets/img/products/upload_<timestamp>.<ext>` et servie via `/static/img/products/upload_<timestamp>.<ext>`.
- Si tu ne fournis pas d’image, ou si le fichier est invalide, le serveur utilise automatiquement l’image par défaut: `/static/img/products/00958058-veste-lazare-gris-moyen_5_1.avif`.
- Côté navigateur, un `onerror` est présent sur les `<img>` (liste et détails): si le chargement échoue, l’image par défaut est utilisée automatiquement.

### Changer l’image par défaut

1. Place ton fichier image dans `assets/img/products/` (ex: `default.webp`).
2. Dans `main.go`, remplace la valeur de `defaultImage` dans le handler POST `/add/create`:

```go
defaultImage := "/static/img/products/default.webp"
```

3. Dans `index.html` et `product.html`, remplace également l’URL dans l’attribut `onerror` pour pointer vers la même image par défaut.



## Règles d’affichage des prix

- Sans réduction: prix en blanc.
- Avec réduction: ancien prix (Price) barré en blanc + nouveau prix (ReducPrice) en rouge à côté.

## Dépannage

- Rien ne change côté styles: force le rechargement (Ctrl+F5). Le template `add.html` référence `style.css?v=2` pour contourner le cache.
- Les images statiques ne s’affichent pas: ouvre l’URL directe, ex. `http://localhost:8080/static/img/products/00958058-veste-lazare-gris-moyen_5_1.avif`. Si 404, vérifie l’emplacement du fichier et le montage du FileServer dans `main.go`.
- Clic sur "Consulter" ne fonctionne pas: assure-toi que la route `/product/` est déclarée avant `/` dans `main.go`.
- Upload refusé: seules les extensions `.png,.jpg,.jpeg,.webp,.avif,.gif` sont acceptées par le serveur.

