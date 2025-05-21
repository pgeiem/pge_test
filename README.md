# QuoteEngine - Moteur de tarification

Le QuoteEngine est un moteur de calcul tarifaire qui permet de définir et d'appliquer des règles de tarification complexes. En analysant ces règles, il génère une table de correspondance entre durées et montants. Cette représentation standardisée permet aux PoS d'obtenir facilement les informations tarifaires nécessaires pour effectuer des calculs de prix ou générer des tickets via un player.


# Grammaire de description des tarifs
 *TODO*


# Format de sortie

Le moteur de tarification génère une structure JSON contenant deux parties principales :
- Un en-tête
- Une liste de segments tarifaires (`table`)

Exemple:
```json
{
    "now": "2023-10-01T08:00:00Z",
    "table": [
        {
            "d": 3600,
            "a": 2.50,
            "l": true,
            "dt": "p"
        },
        // ... autres segments
    ]
}
```

## En-tête

- `now`: Date et heure de référence du début du tarif au format RFC3339
- `expiry`: Date et heure au format RFC3339 définissant le moment où le droit de stationnement associé cesse d'influencer le calcul du prochain quota.

## Table 

### Champs obligatoires

- `d` (Duration) :
  - Durée du segment en secondes
  - Les durées sont cumulatives pour déterminer la position temporelle de chaque segment
  - Exemple : `"d": 3600` représente 1 heure

- `a` (Amount) :
  - Montant dans la devise courante
  - Les montants sont cumulatifs pour obtenir le coût total
  - Exemple : `"a": 2.50` représente 2,50 dans la devise

- `l` (Linear) :
  - `true` : Le tarif est calculé proportionnellement au temps
  - `false` : Le tarif est fixe pour la durée du segment
  - Exemple : `"l": true` pour un tarif horaire, `"l": false` pour un palier 

- `dt` (DurationType) :
  - `p` : Segment payant (Paying)
  - `f` : Segment gratuit (Free)
  - `np` : Segment non payant (Non-paying)
  - `b` : Segment interdit (Banned)

### Champs optionnels

- `n` (Name) :
  - Nom du segment pour identification
  - Exemple : `"n": "tarif_jour"`
  
- `dbg` (Debug) :
  - Liste d'informations de débogage
  - Exemple : `"dbg": ["shift to 0s", "truncate after 1h"]`

- `m` (Meta) :
  - Données additionnelles au format JSON
  - Exemple : `"m": {"color": "red"}`

## Exemples de cas d'usage

### Tarif simple horaire
```json
{
    "now": "2023-10-01T08:00:00Z",
    "table": [{
        "d": 3600,
        "a": 2.50,
        "l": true,
        "dt": "p"
    }]
}
```
Un seul segment linéaire de type payant, d'une durée de 1 heure (3600 secondes), avec un montant de 2,50 (en € / CHF / ...).

### Exemple avec un segment gratuit suivi d'un segment payant
```json
{
    "now": "2023-10-01T08:00:00Z",
    "table": [
        {
            "d": 1800,
            "a": 0,
            "l": false,
            "dt": "f"
        },
        {
            "d": 3600,
            "a": 3.00,
            "l": true,
            "dt": "p"
        }
    ]
}
```
Un segment gratuit de 30 minutes suivi d'un segment payant linéaire d'une heure avec un montant de 3,00.

### Exemples complet avec les champs optionnels
```json
{
	"now": "2025-03-14T15:54:30+01:00",
	"table": [{
		"n": "minamount",
		"dbg": ["Repetition no0", "shift to 0s"],
		"d": 900,
		"a": 0.5,
		"l": false,
		"dt": "p"
	}, {
		"n": "minamount2",
		"dbg": ["Repetition no0", "shift to 15m0s"],
		"d": 900,
		"a": 0.1,
		"l": false,
		"dt": "p"
	}, {
		"n": "30 min linear",
		"dbg": ["shift to 30m0s"],
		"d": 1800,
		"a": 0.8,
		"l": true,
		"dt": "p"
	}, {
		"n": "hourlyrate",
		"dbg": ["shift to 1h0m0s"],
		"d": 3600,
		"a": 1.1,
		"l": true,
		"dt": "p"
	}, {
		"n": "hourlyrate",
		"dbg": ["shift to 2h0m0s", "solve against night", "truncate after 3h5m30s"],
		"d": 3930,
		"a": 1.091667,
		"l": true,
		"dt": "p"
	}, {
		"n": "night",
		"dbg": ["Occurence no0"],
		"d": 50400,
		"a": 0,
		"l": false,
		"dt": "np"
	}, {
		"n": "hourlyrate",
		"dbg": ["shift to 2h0m0s", "solve against night", "truncate split between 3h5m30s and 17h5m30s"],
		"d": 3270,
		"a": 0.908333,
		"l": true,
		"dt": "p"
	}, {
		"n": "3rdhourlyrate",
		"dbg": ["shift to 18h0m0s", "solve against lunch", "truncate after 20h5m30s"],
		"d": 7530,
		"a": 1.045833,
		"l": true,
		"dt": "p"
	}, {
		"n": "lunch",
		"dbg": ["Occurence no0"],
		"d": 7200,
		"a": 0,
		"l": false,
		"dt": "np"
	}, {
		"n": "3rdhourlyrate",
		"dbg": ["shift to 18h0m0s", "solve against lunch", "truncate split between 20h5m30s and 22h5m30s"],
		"d": 3270,
		"a": 0.454167,
		"l": true,
		"dt": "p"
	}, {
		"n": "flat1",
		"dbg": ["Repetition no0", "shift to 23h0m0s"],
		"d": 600,
		"a": 1,
		"l": false,
		"dt": "p"
	}, {
		"n": "flat2",
		"dbg": ["Repetition no0", "shift to 23h10m0s"],
		"d": 600,
		"a": 2,
		"l": false,
		"dt": "p"
	}, {
		"n": "flat3",
		"dbg": ["Repetition no0", "shift to 23h20m0s"],
		"d": 600,
		"a": 8,
		"l": false,
		"dt": "p"
	}, {
		"n": "FPS",
		"dbg": ["Repetition no0", "shift to 23h30m0s"],
		"d": 1800,
		"a": 13,
		"l": false,
		"dt": "p",
		"m": {
			"color": "red"
		}
	}]
}
```

