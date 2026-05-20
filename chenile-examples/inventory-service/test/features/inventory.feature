Feature: Inventory service

  Scenario: Create inventory
    When I POST a REST request to URL "/inventorys" with payload
      """
      {
        "name": "Alice"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Alice"

  Scenario: Reject inventory without name
    When I POST a REST request to URL "/inventorys" with payload
      """
      {
        "name": ""
      }
      """
    Then the http status code is 400
    And success is false
    And the error array size is 2
