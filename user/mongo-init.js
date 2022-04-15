db.getSiblingDB('nabiel').createUser({
  user: 'nabiel',
  pwd: 'nabiel',
  roles: [{
    role: 'readWrite',
    db: 'nabiel'
  }]
})