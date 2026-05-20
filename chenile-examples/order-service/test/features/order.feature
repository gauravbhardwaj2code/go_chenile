Feature: Order service

  Scenario: Create order
    When I POST a REST request to URL "/orders" with payload
      """
      {
        "name": "Alice"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Alice"

  Scenario: Reject order without name
    When I POST a REST request to URL "/orders" with payload
      """
      {
        "name": ""
      }
      """
    Then the http status code is 400
    And success is false
    And the error array size is 2
