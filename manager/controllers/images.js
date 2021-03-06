const express = require('express');
const axios = require('axios');
const models = require('../models');

const router = express.Router();
const url = `${process.env.BUILDER_URL}/images`;

router.get('/', (req, res) => {
  models.Images.findAll().then(images => res.json(images));
});

router.post('/', (req, res, next) => {
  models.Images.create(req.body)
    .then(image => image.get({ plain: true }))
    .then(image => axios.post(url, image))
    .then(resp => resp.data)
    .then(image => res.json(image))
    .catch(err => next(err));
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

router.put('/:id', (req, res, next) => {
  let isBuild = req.body.build;
  if (isBuild) {
    delete req.body.build;
    req.body.path = null;
    req.body.error = null;
  }
  models.Images.findById(req.params.id).then(image => {
    if (image == null) {
      throw new Error('image is null');
    }
    return image.update(req.body);
  }).then(image => isBuild ? axios.post(url, image) : { data: image })
    .then(resp => resp.data)
    .then(image => {
      req.io.emit('image', image);
      res.json(image);
    })
    .catch(err => next(err));
});

router.delete('/:id', (req, res, next) => {
  models.Images.destroy({
    where: {
      id: {
        $eq: req.params.id,
      },
    },
  }).then(() => res.status(204).send());
});

module.exports = router;
