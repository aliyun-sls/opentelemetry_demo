// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
const { context, propagation, trace, metrics } = require('@opentelemetry/api');
const cardValidator = require('simple-card-validator');
const { v4: uuidv4 } = require('uuid');

const { OpenFeature } = require('@openfeature/server-sdk');
const { FlagdProvider } = require('@openfeature/flagd-provider');
const flagProvider = new FlagdProvider();

const logger = require('./logger');
const tracer = trace.getTracer('payment');
const meter = metrics.getMeter('payment');
const transactionsCounter = meter.createCounter('app.payment.transactions');

const LOYALTY_LEVEL = ['platinum', 'gold', 'silver', 'bronze'];

/** Return random element from given array */
function random(arr) {
  const index = Math.floor(Math.random() * arr.length);
  return arr[index];
}

module.exports.charge = async request => {
  const span = tracer.startSpan('charge');

  logger.info({
    amount: request.amount,
    card_last_four: request.creditCard.creditCardNumber.substr(-4),
    card_expiry: `${request.creditCard.creditCardExpirationMonth}/${request.creditCard.creditCardExpirationYear}`
  }, '收到支付请求');

  await OpenFeature.setProviderAndWait(flagProvider);

  const numberVariant = await OpenFeature.getClient().getNumberValue("paymentFailure", 0);

  if (numberVariant > 0) {
    // n% chance to fail with app.loyalty.level=gold
    if (Math.random() < numberVariant) {
      logger.warn({
        loyalty_level: 'gold',
        failure_rate: numberVariant
      }, '支付故障注入被触发');
      
      span.setAttributes({'app.loyalty.level': 'gold' });
      span.end();

      throw new Error('Payment request failed. Invalid token. app.loyalty.level=gold');
    }
  }

  const {
    creditCardNumber: number,
    creditCardExpirationYear: year,
    creditCardExpirationMonth: month
  } = request.creditCard;
  const currentMonth = new Date().getMonth() + 1;
  const currentYear = new Date().getFullYear();
  const lastFourDigits = number.substr(-4);
  const transactionId = uuidv4();

  const card = cardValidator(number);
  const { card_type: cardType, valid } = card.getCardDetails();

  const loyalty_level = random(LOYALTY_LEVEL);

  logger.info({
    transaction_id: transactionId,
    card_type: cardType,
    card_valid: valid,
    loyalty_level: loyalty_level
  }, '开始处理支付');

  span.setAttributes({
    'app.payment.card_type': cardType,
    'app.payment.card_valid': valid,
    'app.loyalty.level': loyalty_level
  });

  if (!valid) {
    logger.error({
      transaction_id: transactionId,
      card_type: cardType,
      card_last_four: lastFourDigits
    }, '信用卡信息无效');
    throw new Error('Credit card info is invalid.');
  }

  if (!['visa', 'mastercard'].includes(cardType)) {
    logger.error({
      transaction_id: transactionId,
      card_type: cardType
    }, '不支持的信用卡类型');
    throw new Error(`Sorry, we cannot process ${cardType} credit cards. Only VISA or MasterCard is accepted.`);
  }

  if ((currentYear * 12 + currentMonth) > (year * 12 + month)) {
    logger.error({
      transaction_id: transactionId,
      card_last_four: lastFourDigits,
      expiry_date: `${month}/${year}`
    }, '信用卡已过期');
    throw new Error(`The credit card (ending ${lastFourDigits}) expired on ${month}/${year}.`);
  }

  // Check baggage for synthetic_request=true, and add charged attribute accordingly
  const baggage = propagation.getBaggage(context.active());
  if (baggage && baggage.getEntry('synthetic_request') && baggage.getEntry('synthetic_request').value === 'true') {
    span.setAttribute('app.payment.charged', false);
    logger.info({
      transaction_id: transactionId,
      charged: false
    }, '模拟请求，不实际扣款');
  } else {
    span.setAttribute('app.payment.charged', true);
    logger.info({
      transaction_id: transactionId,
      charged: true
    }, '实际扣款');
  }

  const { units, nanos, currencyCode } = request.amount;
  logger.info({
    transaction_id: transactionId,
    card_type: cardType,
    card_last_four: lastFourDigits,
    amount: { units, nanos, currencyCode },
    loyalty_level
  }, '支付处理完成');
  
  transactionsCounter.add(1, { 'app.payment.currency': currencyCode });
  span.end();

  return { transactionId };
};
