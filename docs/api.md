# API V1 Documentation

## General Options

### Field Limiting
All endpoints of this API offer field limiting by adding a `fields` parameter to the request. The response JSON will then only contain the fields provided as values for this parameter, all other fields will be omitted. Non-existent field names will be ignored, if only non-existent fields are provided, the JSON will empty. The values of the `fields` parameter need to be separated by commata. Example: `v1/pokemon/1?fields=name,classification`

### Sorting
All lists of resources offer sorting by id or name of the resources with the query parameter `sort`.
* Options are: `id_asc`, `id_desc`, `name_asc`, `name_desc`
* Only the first value provided is used for sorting.

### Pagination
TODO: Add definition for limiting results and using a offset.

## General Types
### NamedResource
This type represents a single API resources and is used in lists of resources as a short representation.
| Name        | Description                                                              | Type          |
| ----------- | ------------------------------------------------------------------------ | ------------- |
| name        | The name of the resource.                                                | String        |
| url         | The URL of this API that offers detailed information about the resource. | String        |

## Abilities
### `GET` **/v1/abilities**
Returns a list of all abilities.
```json
{
  "count": <number of abilities>,
  "results": [
    {
      "name": "<ability-name>",
      "url": "<instance-url>/abilities/<ability-id>"
    }
  ]
}
```
#### **AbilityList**
| Name        | Description                                                | Type                   |
| ----------- | ---------------------------------------------------------- | ---------------------- |
| count       | Total number of ability resources available from this API. | Integer                |
| results     | A list of named ability resources.                         | Array\<NamedResource\> |

### `GET` **/v1/abilities/_\<id or name\>_**
Returns data about a single ability.
```json
{
  "id": <ability-id>,
  "name": "<ability-name>",
  "description": "<ability-description>",
  "pokemon": [
    {
      "name": "<pokemon-name>",
      "url": "<instance-url>/pokemon/<pokemon-id>"
    }
  ]
}
```
#### **Ability**
| Name        | Description                                                | Type                   |
| ----------- | ---------------------------------------------------------- | ---------------------- |
| id          |                                                            | Integer                |
| name        |                                                            | String                 |
| description |                                                            | String                 |
| pokemon     |                                                            | Array\<NamedResource\> |

## Camps
### `GET` **/v1/camps**
Returns a list of all camps.
```json
{
  "count": <number of camps>,
  "results": [
    {
      "name": "<camp-name>",
      "url": "<instance-url>/camps/<camp-id>"
    }
  ]
}
```
#### **CampList**
| Name        | Description                                             | Type                   |
| ----------- | ------------------------------------------------------- | ---------------------- |
| count       | Total number of camp resources available from this API. | Integer                |
| results     | A list of named camp resources.                         | Array\<NamedResource\> |

### `GET` **/v1/camps/_\<id or name\>_**
Returns data about a single camp.
```json
{
  "id": <camp-id>,
  "name": "<camp-name>",
  "description": "<camp-description>",
  "unlockType": "<unlock-type>",
  "cost": <cost>,
  "pokemon": [
    {
      "name": "<pokemon-name>",
      "url": "<instance-url>/pokemon/<pokemon-id>"
    }
  ]
}
```
#### **Camp**
| Name        | Description                                                | Type                   |
| ----------- | ---------------------------------------------------------- | ---------------------- |
| id          |                                                            | Integer                |
| name        |                                                            | String                 |
| description |                                                            | String                 |
| unlockType  |                                                            | String                 |
| cost        |                                                            | Integer                |
| pokemon     |                                                            | Array\<NamedResource\> |

## Dungeons
### `GET` **/v1/dungeons**
Returns a list of all dungeons.
```json
{
  "count": <number of dungeons>,
  "results": [
    {
      "name": "<dungeon-name>",
      "url": "<instance-url>/dungeons/<dungeon-id>"
    }
  ]
}
```
#### **DungeonList**
| Name        | Description                                                | Type                   |
| ----------- | ---------------------------------------------------------- | ---------------------- |
| count       | Total number of dungeon resources available from this API. | Integer                |
| results     | A list of named dungeon resources.                         | Array\<NamedResource\> |


### `GET` **/v1/dungeons/_\<id or name\>_**
Returns data about a single dungeon.
```json
{
  "id": <dungeon-id>,
  "name": "<dungeon-name>",
  "levels": <levels>,
  "startLevel": <start-level>,
  "teamSize": <team-size>,
  "itemsAllowed": <items-allowed>,
  "pokemonJoining": <pokemon-joining>,
  "mapVisible": <map-visible>,
  "pokemon": [
    {
      "pokemon": {
        "name": "<pokemon-name>",
        "url": "<instance-url>/pokemon/<pokemon-id>"
      },
      "isSuper": <super_pokemon>
    }

  ]
}
```
#### **Dungeon**
| Name           | Description                                                | Type                    |
| -------------- | ---------------------------------------------------------- | ----------------------- |
| id             |                                                            | Integer                 |
| name           |                                                            | String                  |
| levels         |                                                            | Integer                 |
| startLevel     |                                                            | Integer                 |
| teamSize       |                                                            | Integer                 |
| itemsAllowed   |                                                            | Boolean                 |
| pokemonJoining |                                                            | Boolean                 |
| mapVisible     |                                                            | Boolean                 |
| pokemon        |                                                            | Array\<DungeonPokemon\> |

#### **DungeonPokemon**
| Name        | Description                                                | Type              |
| ----------- | ---------------------------------------------------------- | ------------------|
| pokemon     |                                                            | \<NamedResource\> |
| isSuper     |                                                            | Boolean           |

## Moves
### `GET` **/v1/moves**
Returns a list of all moves.
```json
{
  "count": <number of moves>,
  "results": [
    {
      "name": "<move-name>",
      "url": "<instance-url>/moves/<move-id>"
    }
  ]
}
```
#### **MoveList**
| Name        | Description                                             | Type                   |
| ----------- | ------------------------------------------------------- | ---------------------- |
| count       | Total number of move resources available from this API. | Integer                |
| results     | A list of named move resources.                         | Array\<NamedResource\> |

### `GET` **/v1/moves/_\<id or name\>_**
Returns data about a single move.
```json
{
  "id": <move-id>,
  "name": "<move-name>",
  "category": <move-category>,
  "range": <move-range>,
  "target": <move-target>,
  "initialPP": <initial_pp>,
  "initialPower": <initial_power>,
  "accuracy": <accuracy>,
  "description": "<move-description>",
  "type": {
    "name": "<type-name>",
    "url": "<instance-url>/types/<type-id>"
  },
  "pokemon": [
    {
      "pokemon": {
        "name": "<pokemon-name>",
        "url": "<instance-url>/pokemon/<pokemon-id>"
      },
      "method": "<learn-type>",
      "level": <level>,
      "cost": <cost>
    }
  ]
}
```
#### **Move**
| Name         | Description                                                | Type                 |
| ------------ | ---------------------------------------------------------- | -------------------- |
| id           |                                                            | Integer              |
| name         |                                                            | String               |
| category     |                                                            | String               |
| range        |                                                            | String               |
| target       |                                                            | String               |
| initialPP    |                                                            | Integer              |
| initialPower |                                                            | Integer              |
| accuracy     |                                                            | Integer              |
| description  |                                                            | String               |
| type         |                                                            | NamedResource        |
| pokemon      |                                                            | Array\<MovePokemon\> |

#### **MovePokemon**
| Name        | Description                                                | Type              |
| ----------- | ---------------------------------------------------------- | ----------------- |
| pokemon     |                                                            | \<NamedResource\> |
| method      |                                                            | String            |
| level       |                                                            | Integer           |
| cost        |                                                            | Integer           |

## Pokemon
### `GET` **/v1/pokemon**
Returns a list of all Pokemon.
```json
{
  "count": <number of pokemon>,
  "results": [
    {
      "name": "<pokemon-name>",
      "url": "<instance-url>/pokemon/<pokemon-id>"
    }
  ]
}
```
#### **PokemonList**
| Name        | Description                                                | Type                   |
| ----------- | ---------------------------------------------------------- | ---------------------- |
| count       | Total number of pokemon resources available from this API. | Integer                |
| results     | A list of named pokemon resources.                         | Array\<NamedResource\> |


#### Filtering
This resource list offers filters as query parameters:
* By type: `type=<type-name or type-id>`
* By move: `move=<move-name or move-id>`

To filter by multiple values, simply separate the values with commata. Example: `/v1/pokemon?type=fire,flight`

### `GET` **/v1/pokemon/_\<id or name\>_:**
Returns data about a single pokemon.
```json
{
  "id": <dex-id>,
  "name": "<pokemon-name>",
  "classification": "<classification>",
  "evolutionStage": <stage-number>,
  "evolveCondition": "<evolve-condition>",
  "evolveLevel": <evolve-level>,
  "evolveCrystals": <number of crystals>,
  "camp": {
    "name": "<camp-name>",
    "url": "<instance-url>/camps/<camp-id>"
  },
  "abilities": [
    {
      "name": "<ability-name>",
      "url": "<instance-url>/abilities/<ability-id>"
    }
  ],
  "dungeons": [
    {
      "dungeon": {
        "name": "<dungeon-name>",
        "url": "<instance-url>/dungeons/<dungeon-id>"
      },
      "isSuper": <super_pokemon>
    }
  ],
  "moves": [
    {
      "move": {
        "name": "<move-name>",
        "url": "<instance-url>/moves/<move-id>"
      },
      "method": "<learn-type>",
      "level": <level>,
      "cost": <cost>
    }
  ],
  "types": [
    {
      "name": "<type-name>",
      "url": "<instance-url>/types/<type-id>"
    }
  ]
}
```

#### **Pokemon**
| Name            | Description                                                | Type                    |
| --------------- | ---------------------------------------------------------- | ----------------------- |
| id              |                                                            | Integer                 |
| name            |                                                            | String                  |
| classification  |                                                            | String                  |
| evolutionStage  |                                                            | Integer                 |
| evolveCondition |                                                            | String                  |
| evolveLevel     |                                                            | Integer                 |
| evolveCrystals  |                                                            | Integer                 |
| camp            |                                                            | NamedResource           |
| abilities       |                                                            | Array\<NamedResource\>  | 
| dungeons        |                                                            | Array\<PokemonDungeon\> |
| moves           |                                                            | Array\<PokemonMove\>    |
| types           |                                                            | Array\<NamedResource\>  |


#### **PokemonDungeon**
| Name        | Description                                                | Type              |
| ----------- | ---------------------------------------------------------- | ------------------|
| dungeon     |                                                            | \<NamedResource\> |
| isSuper     |                                                            | Boolean           |

#### **PokemonMove**
| Name        | Description                                                | Type              |
| ----------- | ---------------------------------------------------------- | ----------------- |
| move        |                                                            | \<NamedResource\> |
| method      |                                                            | String            |
| level       |                                                            | Integer           |
| cost        |                                                            | Integer           |

## Types
### `GET` **/v1/types**
Returns a list of all types.
```json
{
  "count": <number of types>,
  "results": [
    {
      "name": "<type-name>",
      "url": "<instance-url>/types/<type-id>"
    }
  ]
}
```
#### **TypeList**
| Name        | Description                                             | Type                   |
| ----------- | ------------------------------------------------------- | ---------------------- |
| count       | Total number of type resources available from this API. | Integer                |
| results     | A list of named type resources.                         | Array\<NamedResource\> |

### `GET` **/v1/types/_\<id or name\>_**
Returns data about a single type.
```json
{
  "id": <type-id>,
  "name": "<type-name>",
  "interactions": [
    {
      "defender": {
        "name": "<type-name>",
        "url": "<instance-url>/types/<type-id>"
      },
      "interaction": <effectiveness>
    }
  ]
}
```
#### **Type**
| Name         | Description                                                | Type                     |
| ------------ | ---------------------------------------------------------- | ------------------------ |
| id           |                                                            | Integer                  |
| name         |                                                            | String                   |
| interactions |                                                            | Array\<TypeInteraction\> |

#### **TypeInteraction**
| Name        | Description                                                | Type              |
| ----------- | ---------------------------------------------------------- | ------------------|
| defender    |                                                            | \<NamedResource\> |
| interaction |                                                            | String            |