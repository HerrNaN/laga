# Manual

Run with admin access

## Acts on

| Action | Step      | Resource                                       |
| ------ | --------- | ---------------------------------------------- |
| Create | Bootstrap | Pulumi state bucket (which only admin can use) |

# Bootstrap

Run using pulumi with admin access

## Prereq.

- Bootstrap - Pulumi state bucket

## Acts on

| Action | Step      | Resource                          |
| ------ | --------- | --------------------------------- |
| Modify | Bootstrap | Pulumi state bucket               |
| Create | Deploy    | Pulumi state bucket               |
| Create | Deploy    | Application                       |
| Create | Deploy    | Minimal application access policy |
| Create | Deploy    | API Key                           |

# Deploy

Run using pulumi with Deploy API Key

## Prerequisites

- Deploy - Pulumi state bucket
- Application
- API Key

## Acts on

| Action | Step   | Resource             |
| ------ | ------ | -------------------- |
| Modify | Deploy | Pulumi state bucket  |
| Create | Deploy | Serverless Container |
