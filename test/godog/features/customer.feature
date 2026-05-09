Feature: Customer service

  Scenario: Create customer
    When I POST a REST request to URL "/customers" with payload
      """
      {
        "name": "Alice"
      }
      """
    Then the http status code is 200
    And success is true
    And the REST response key "name" is "Alice"

