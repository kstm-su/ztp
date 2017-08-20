const express = require('express');
const router = express.Router();

router.get('/', (req, res) => {
  req.nodes.findAll({
    include: [req.images],
  }).then(nodes => res.json(nodes));
});

router.post('/', (req, res, next) => {
  req.nodes.create(req.body).then(node => {
    res.json(node.get({
      plain: true,
    }));
  }).catch(err => next(err));
});

router.get('/:id', (req, res, next) => {
  req.nodes.findById(req.params.id, {
    include: [req.images],
  }).then(node => {
    if (node == null) {
      return next();
    }
    res.json(node);
  });
});

router.put('/:id', (req, res, next) => {
  req.nodes.findById(req.params.id).then(node => {
    if (node == null) {
      return next();
    }
    node.update(req.body).then(node => res.json(node)).catch(err => next(err));
  });
});

module.exports = router;
