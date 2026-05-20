Feature: Slot service allocation

  Scenario: Allocate runner that satisfies hard and soft constraints
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Asha",
        "skills": ["cook"],
        "attributes": {
          "diet": "veg",
          "gender": "female"
        },
        "slots": [
          {"date": "2026-06-01", "start": "09:00", "end": "11:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Ravi",
        "skills": ["cook"],
        "attributes": {
          "diet": "veg",
          "gender": "male"
        },
        "slots": [
          {"date": "2026-06-01", "start": "09:00", "end": "11:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/allocations" with payload
      """
      {
        "requestId": "req-1",
        "skill": "cook",
        "slot": {"date": "2026-06-01", "start": "09:00", "end": "11:00"},
        "constraints": [
          {"key": "diet", "value": "veg", "type": "hard"},
          {"key": "gender", "value": "female", "type": "soft"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "runnerName" is "Asha"
    And the REST response key "softScore" is "1"

  Scenario: Reject allocation when hard constraints cannot be satisfied
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Kiran",
        "skills": ["cook"],
        "attributes": {
          "diet": "non-veg"
        },
        "slots": [
          {"date": "2026-06-02", "start": "09:00", "end": "11:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/allocations" with payload
      """
      {
        "requestId": "req-2",
        "skill": "cook",
        "slot": {"date": "2026-06-02", "start": "09:00", "end": "11:00"},
        "constraints": [
          {"key": "diet", "value": "veg", "type": "hard"}
        ]
      }
      """
    Then the http status code is 404
    And success is false
    And the error array size is 1

  Scenario: Prefer the runner with the best soft constraint score
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Mohan",
        "skills": ["cook"],
        "attributes": {
          "diet": "veg",
          "gender": "male"
        },
        "slots": [
          {"date": "2026-06-04", "start": "08:00", "end": "10:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Leela",
        "skills": ["cook"],
        "attributes": {
          "diet": "veg",
          "gender": "female"
        },
        "slots": [
          {"date": "2026-06-04", "start": "08:00", "end": "10:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/allocations" with payload
      """
      {
        "requestId": "req-soft",
        "skill": "cook",
        "slot": {"date": "2026-06-04", "start": "08:00", "end": "10:00"},
        "constraints": [
          {"key": "diet", "value": "veg", "type": "soft"},
          {"key": "gender", "value": "female", "type": "soft"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "runnerName" is "Leela"
    And the REST response key "softScore" is "2"

  Scenario: Prevent double booking of the same runner slot
    When I POST a REST request to URL "/runners" with payload
      """
      {
        "name": "Meena",
        "skills": ["house-help"],
        "attributes": {
          "gender": "female"
        },
        "slots": [
          {"date": "2026-06-03", "start": "15:00", "end": "17:00"}
        ]
      }
      """
    Then the http status code is 200
    And success is true
    When I POST a REST request to URL "/allocations" with payload
      """
      {
        "requestId": "req-3",
        "skill": "house-help",
        "slot": {"date": "2026-06-03", "start": "15:00", "end": "17:00"}
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "runnerName" is "Meena"
    When I POST a REST request to URL "/allocations" with payload
      """
      {
        "requestId": "req-4",
        "skill": "house-help",
        "slot": {"date": "2026-06-03", "start": "15:00", "end": "17:00"}
      }
      """
    Then the http status code is 404
    And success is false
