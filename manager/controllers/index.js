const express = require('express');
const router = express.Router();

router.use('/images', require('./images'));
router.use('/nodes', require('./nodes'));

router.use((req, res, next) => {
  next({
    message: `'${req.originalUrl}' does not exist`,
    status: 404,
  });
});

router.use((err, req, res, next) => {
  if (err.status == null) {
    err.status = 500;
  }
  console.error(err);
  res.status(err.status).json({
    error: err,
  });
});

module.exports = router;
