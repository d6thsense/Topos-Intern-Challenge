#curl -d '{"Borough":"Test", "Status":"OK", "ConstructionYear":346, "Height":234, "Type":"Nice"}' -H "Content-Type: application/json" -X POST http://localhost:4018/addBuilding >> TestRes.json

#curl http://localhost:4018/statHeightByType >> TestRes.json

#curl http://localhost:4018/getBuildings >> TestRes.json

curl -X "DELETE" http://localhost:4018/removeBuilding/5ca41979ba6d3392e73be3d3
