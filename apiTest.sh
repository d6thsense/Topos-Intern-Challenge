curl -d '{"Borough":"Test1", "Status":"OK", "ConstructionYear":345, "Height":234, "Type":"Nice"}' -H "Content-Type: application/json" -X POST http://localhost:4018/addBuilding >> TestRes.json

curl http://localhost:4018/statHeightByType >> TestRes.json

curl http://localhost:4018/buildings >> TestRes.json

