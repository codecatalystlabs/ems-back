package bootstrap

import (
	"dispatch/internal/platform/config"

	"github.com/IBM/sarama"
)

func NewKafkaSyncProducer(cfg config.KafkaConfig) (sarama.SyncProducer, error) {
	scfg := sarama.NewConfig()
	scfg.ClientID = cfg.ClientID
	scfg.Version = sarama.V3_6_0_0
	scfg.Producer.Return.Successes = true
	scfg.Producer.RequiredAcks = sarama.WaitForAll
	return sarama.NewSyncProducer(cfg.Brokers, scfg)
}

func NewKafkaConsumerGroup(cfg config.KafkaConfig) (sarama.ConsumerGroup, error) {
	scfg := sarama.NewConfig()
	scfg.ClientID = cfg.ClientID
	scfg.Version = sarama.V3_6_0_0
	scfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	return sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, scfg)
}
