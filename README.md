To run the Auctioneer

    make
    make run

To run the bidders

        cd cmd/bidder
        go build
        go run main.go -name=<bidder name> -port=<port> -delay="bidders delay time"
    

Endpoints:
 
    To get all the bidders - 
            ```
                curl -X GET \
                http://localhost:5000/v1/bidder/all \
                -H 'cache-control: no-cache' \
            ```

    To trigger and auction - 
            ```
                curl -X POST \
                http://localhost:5000/v1/auction \
                -H 'cache-control: no-cache' \
                -H 'content-type: application/json' \
                -d '{
                    "auction_id": "Sample"
                    }'
            ```


    To register a bidder against an auctioneer - 
        ```
            curl -X POST \
            http://localhost:5000/v1/bidder/register \
            -H 'cache-control: no-cache' \
            -H 'content-type: application/json' \
            -d '{
                "name": "Test",
                "delay": 50
                }'
        ```

    Bidding - 
            ```
                curl -X POST \
                http://localhost:5000/v1/auction \
                -H 'cache-control: no-cache' \
                -H 'content-type: application/json' \
                -H 'postman-token: 3347d70f-6240-36ed-14ae-633773bb4c8e' \
                -d '{
                    "auction_id": "some_random_id"
                }'
            ```

    

    