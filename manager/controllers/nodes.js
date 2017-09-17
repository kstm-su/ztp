const express = require('express');
const models = require('../models');

const router = express.Router();

router.get('/', (req, res) => {
  models.Nodes.findAll({
    include: [models.Images],
  }).then(nodes => res.json(nodes));
});

router.post('/', (req, res, next) => {
  models.Nodes.create(req.body).then(node => {
    res.json(node.get({
      plain: true,
    }));
  }).catch(err => next(err));
});

router.get('/:id', (req, res, next) => {
  models.Nodes.findById(req.params.id, {
    include: [models.Images],
  }).then(node => {
    if (node == null) {
      return next();
    }
    res.json(node);
  });
});

router.put('/:id', (req, res, next) => {
  models.Nodes.findById(req.params.id).then(node => {
    if (node == null) {
      return next();
    }
    node.update(req.body)
      .then(node => {
        req.io.emit('node', node);
        res.json(node);
      })
      .catch(err => next(err));
  });
});

router.delete('/:id', (req, res, next) => {
  models.Nodes.destroy({
    where: {
      id: {
        $eq: req.params.id,
      },
    },
  }).then(() => res.status(204).send());
});

module.exports = router;
