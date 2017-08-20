const express = require('express');
const router = express.Router();

router.get('/', (req, res) => {
  req.images.findAll().then(images => res.json(images));
});

router.post('/', (req, res, next) => {
  req.images.create(req.body).then(image => {
    res.json(image.get({
      plain: true,
    }));
  }).catch(err => next(err));
});

router.get('/:id', (req, res, next) => {
  req.images.findById(req.params.id, { 
    include: [req.nodes]
  }).then(image => {
    if (image == null) {
      return next();
    }
    res.json(image);
  });
});

router.put('/:id', (req, res) => {
  req.images.findById(req.params.id).then(image => {
    if (image == null) {
      return next();
    }
    image.update(req.body).then(image => res.json(image)).catch(err => next(err));
  });
});

module.exports = router;
