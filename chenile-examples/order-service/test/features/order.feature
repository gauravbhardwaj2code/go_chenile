Feature: Order service

  Scenario: Create order
    When I POST a REST request to URL "/orders" with payload
      """
      {
        "name": "Order 1"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Order 1"

