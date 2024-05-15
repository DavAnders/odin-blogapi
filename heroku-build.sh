#!/bin/sh
cd frontend
npm install
npm run build
cd ..
mv frontend/dist/* backend/public/
