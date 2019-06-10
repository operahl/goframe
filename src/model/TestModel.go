package model
/**

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 AUTO_INCREMENT=3 ;

INSERT INTO `user` (`id`, `name`) VALUES
(1, 'jack'),
(2, 'merry');

 */

type TestUser struct {

	Id  int `json:"id" xorm:"id"`

	Name string	`json:"name" xorm:"name"`

}