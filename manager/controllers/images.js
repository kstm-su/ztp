const express = require('express');
const models = require('../models');

const router = express.Router();

router.get('/', (req, res) => {
  models.Images.findAll().then(images => res.json(images));
});

router.post('/', (req, res, next) => {
  models.Images.create(req.body).then(image => {
    res.json(image.get({
      plain: true,
    }));
  }).catch(err => next(err));
});

router.get('/:id', (req, res, next) => {
  models.Images.findById(req.params.id, { 
    include: [models.Nodes]
  }).then(image => {
    if (image == null) {
      return next();
    }
    res.json(image);
  });
});

router.put('/:id', (req, res) => {
  models.Images.findById(req.params.id).then(image => {
    if (image == null) {
      return next();
    }
    image.update(req.body).then(image => res.json(image)).catch(err => next(err));
  });
});

module.exports = router;
