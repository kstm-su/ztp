const Sequelize = require('sequelize');
const sequelize = new Sequelize(process.env.DATABASE_URL || 'sqlite://db.sqlite3');
const axios = require('axios');
const fs = require('fs');

const url = `${process.env.STORAGE_URL}/images`;

const createStandby = (config) => {
  require('./models').Images.findAll().then(images => {
    console.log(images);
    if(images.length == 0){
      console.log("Create image stand by dhcp");
      console.log(config);
      require('./models').Images.create(config)
        .then(image => image.get({ plain: true }))
        .then(image => axios.post(url, image))
        .catch(err => {
          console.log(err);
        });
    }
  });
};

standByConfig = fs.readFileSync('./standby.conf', 'utf8');
standByConfig = standByConfig.replace("${start_ip}", `${process.env.DHCP_START_IP_ADDR}`)
	.replace("${lease_range}", `${process.env.DHCP_LEASE_RANGE}`)
	.replace("${manager_addr}", `${process.env.DHCP_SERVER_IP_ADDR}:${process.env.MANAGER_PORT}`);

var standByRequest = {
	name: 'Image for standby',
	config: standByConfig,
	description: 'An image to registrate mac address and stand by user action',
	size: 512,
	build: true
};

sequelize.authenticate().then(() => {
  console.log('success to connect db');
  require('./models').migration();
}).then(() => {
  console.log("creating stand by image");
  createStandby(standByRequest);
}).catch(err => {
  console.error('unable to connect db:', err);
  process.exit();
});

module.exports = sequelize;
