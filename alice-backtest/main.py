import argparse
import json
import urllib.request
from dataclasses import dataclass
import marshmallow_dataclass

import backtrader as bt
import datetime

import pandas as pd


@dataclass
class Candle:
    # 通貨名
    instrument: str
    # 足種
    granularity: str
    # 始値
    open_price: str
    # 高値
    high_price: str
    # 安値
    low_price: str
    # 終値
    close_price: str
    # 時間
    time: str


class BuySellArrows(bt.observers.BuySell):
    plotlines = dict(buy=dict(marker='$\u21E7$', markersize=12.0),
                     sell=dict(marker='$\u21E9$', markersize=12.0))


class AliceStrategy(bt.Strategy):
    params = dict(
        instrument="",
        granularity=""
    )

    def log(self, txt, dt=None):
        dt = dt or self.datas[0].datetime.date(0)
        print('%s, %s' % (dt.isoformat(), txt))

    def __init__(self):
        # data0=instrument
        self.instrument = self.params.instrument
        # data1=granularity
        self.granularity = self.params.granularity

        # data0=12hours , data1=day となっている。
        self.open_12 = self.data0.open
        self.high_12 = self.data0.high
        self.low_12 = self.data0.low
        self.close_12 = self.data0.close
        self.time_12 = self.data0.datetime

        self.open_24 = self.data1.open
        self.high_24 = self.data1.high
        self.low_24 = self.data1.low
        self.close_24 = self.data1.close
        self.time_24 = self.data1.datetime

    def next(self):
        # --- Candle Data
        candle_12 = Candle(self.instrument,
                           self.granularity,
                           str(self.open_12[0]),
                           str(self.high_12[0]),
                           str(self.low_12[0]),
                           str(self.close_12[0]),
                           self.time_12.datetime(0)
                           )
        candle_24 = Candle(self.instrument,
                           self.granularity,
                           str(self.open_24[0]),
                           str(self.high_24[0]),
                           str(self.low_24[0]),
                           str(self.close_24[0]),
                           self.time_24.datetime(0)
                           )

        # --- Parameter Area Start ---
        # URL -> セットアップ / トレード計画
        url_setup = 'http://localhost:7070/captain.america/setup'
        url_trade_plan = 'http://localhost:7070/captain.america/trade.plan'

        # Request Header -> 共通
        headers = {'Content-Type': 'application/json'}

        # Request Body -> 12時間足 / 24時間足
        param_candle_12 = marshmallow_dataclass.class_schema(Candle)().dump(candle_12)
        param_candle_24 = marshmallow_dataclass.class_schema(Candle)().dump(candle_24)
        # --- Parameter Area End ---

        # --- Back Test Area Start ---
        # ----- Handle Position -----
        pos_size = self.getposition(self.data0).size
        if pos_size < 0:
            if self.datas[1].close[-1] < float(candle_24.close_price):
                self.close()
        elif pos_size > 0:
            if self.datas[1].close[-1] > float(candle_24.close_price):
                self.close()

        # ----- Start 12-candle Check Trade Plan -----
        req_trade_plan_12 = urllib.request.Request(url_trade_plan, json.dumps(param_candle_12).encode(), headers)
        with urllib.request.urlopen(req_trade_plan_12) as res:
            body = res.read()
        trade_plan_12 = json.loads(body)

        # トレード実行
        if trade_plan_12.get("is_order"):
            if trade_plan_12.get("buy_sell") == 'BUY':
                if pos_size >= 0:
                    self.buy(size=1000, exectype=bt.Order.StopLimit, plimit=(float(candle_24.close_price) * 1.003))
                    print('BUY CREATE - 12,', candle_24.time, candle_24.close_price, pos_size)
            elif trade_plan_12.get("buy_sell") == 'SELL':
                if pos_size <= 0:
                    self.sell(size=1000, exectype=bt.Order.StopLimit, plimit=(float(candle_24.close_price) * 1.003))
                    print('SELL CREATE - 12,', candle_24.time, candle_24.close_price, pos_size)
        # ----- End 12-candle Check Trade Plan -----

        # ----- Start 24-candle Check Set Up -----
        req_setup_24 = urllib.request.Request(url_setup, json.dumps(param_candle_24).encode(), headers)
        urllib.request.urlopen(req_setup_24)
        # ----- End 24-candle Check Set Up -----

        # ----- Start 24-candle Check Trade Plan -----
        req_trade_plan_24 = urllib.request.Request(url_trade_plan, json.dumps(param_candle_24).encode(), headers)
        with urllib.request.urlopen(req_trade_plan_24) as res:
            body = res.read()
        trade_plan_24 = json.loads(body)

        # トレード実行
        if trade_plan_24.get("is_order"):
            if trade_plan_24.get("buy_sell") == 'BUY':
                if pos_size >= 0:
                    self.buy(size=1000, exectype=bt.Order.StopLimit, plimit=(float(candle_24.close_price) * 1.003))
                    print('BUY CREATE - 24,', candle_24.time, candle_24.close_price, pos_size)
            elif trade_plan_24.get("buy_sell") == 'SELL':
                if pos_size <= 0:
                    self.sell(size=1000, exectype=bt.Order.StopLimit, plimit=(float(candle_24.close_price) * 1.003))
                    print('SELL CREATE - 24,', candle_24.time, candle_24.close_price, pos_size)
        # ----- End 24-candle Check Trade Plan -----
        # --- Back Test Area End ---


def run_start(args=None):
    args = parse_args(args)
    cerebro = bt.Cerebro(stdstats=False)
    cerebro.addstrategy(AliceStrategy, instrument=args.instrument, granularity=args.granularity)
    # 複数の足種を利用する際は、小さい足種からデータを読み込む必要がある。
    # 12時間足を読み込む。
    data_12hour = bt.feeds.YahooFinanceCSVData(
        dataname='./data/' + args.instrument + '-H12.csv',
        fromdate=datetime.datetime(2016, 1, 1),
        todate=datetime.datetime(2020, 7, 1),
        timeframe=bt.TimeFrame.Minutes,
        compression=720,
        round=False,
        reverse=False)
    cerebro.adddata(data_12hour)

    # 日足を読み込む。
    data_day = bt.feeds.YahooFinanceCSVData(
        dataname='./data/' + args.instrument + '-D.csv',
        fromdate=datetime.datetime(2016, 1, 1),
        todate=datetime.datetime(2020, 7, 1),
        timeframe=bt.TimeFrame.Days,
        round=False,
        reverse=False)
    cerebro.adddata(data_day)

    cerebro.broker.setcash(1000000.0)
    cerebro.addobserver(bt.obs.Broker)
    cerebro.addobserver(bt.obs.Trades)
    cerebro.addobserver(bt.obs.BuySell)
    cerebro.broker.set_coc(True)

    import backtrader.analyzers as btanalyzers
    cerebro.addanalyzer(btanalyzers.DrawDown, _name='DrawDown')
    cerebro.addanalyzer(btanalyzers.TradeAnalyzer, _name='TradeAnalyzer')
    cerebro.addanalyzer(btanalyzers.SQN, _name='SQN')

    print('Starting Cash Value: %.2f' % cerebro.broker.getcash())
    thestrats = cerebro.run(stdstats=False)

    # 実行結果
    startcash = cerebro.broker.getvalue()
    thestrat = thestrats[0]
    print('Max. DrawDown[%]      :{0:.2f}'.format(  # 最大下落率
        thestrat.analyzers.DrawDown.get_analysis().max.drawdown
    ))
    print('Trades                :{0}'.format(  # 取引回数
        thestrat.analyzers.TradeAnalyzer.get_analysis().total.closed
    ))
    print('Win Trades            :{0}'.format(  # 勝ち取引回数
        thestrat.analyzers.TradeAnalyzer.get_analysis().won.total
    ))
    print('Win Avg Price         :{0}'.format(  # 勝ち取引平均利益
        thestrat.analyzers.TradeAnalyzer.get_analysis().won.pnl.average
    ))
    print('Lose Avg Price        :{0}'.format(  # 負け取引平均損失
        thestrat.analyzers.TradeAnalyzer.get_analysis().lost.pnl.average
    ))
    win_rate = (  # 勝率
            thestrat.analyzers.TradeAnalyzer.get_analysis().won.total
            / thestrat.analyzers.TradeAnalyzer.get_analysis().total.closed
    )
    lost_rate = (  # 敗率
            thestrat.analyzers.TradeAnalyzer.get_analysis().lost.total
            / thestrat.analyzers.TradeAnalyzer.get_analysis().total.closed
    )
    print('Win Rate[%]           :{0:.2f}'.format(  # 勝率=勝ち取引回数÷全取引回数×100
        win_rate * 100
    ))
    print('Lose Rate[%]          :{0:.2f}'.format(  # 敗率=負け取引回数÷全取引回数×100
        lost_rate * 100
    ))
    print('Best Trade[%]         :{0:.2f}'.format(  # 1回の取引での利益の最大値÷所持金×100
        thestrat.analyzers.TradeAnalyzer.get_analysis().won.pnl.max / startcash * 100
    ))
    print('Worst Trade[%]        :{0:.2f}'.format(  # 1回の取引での損失の最大値÷所持金×100
        thestrat.analyzers.TradeAnalyzer.get_analysis().lost.pnl.max / startcash * 100
    ))
    print('Avg. Trade[%]         :{0:.2f}'.format(  # 損益の平均値÷所持金×100
        thestrat.analyzers.TradeAnalyzer.get_analysis().pnl.net.average / startcash * 100
    ))
    print('Max. Trade Duration   :{0}'.format(  # 1回の取引での最長期間
        pd.Timedelta(minutes=thestrat.analyzers.TradeAnalyzer.get_analysis().len.max)
    ))
    print('Avg. Trade Duration   :{0}'.format(  # 1回の取引での平均期間
        pd.Timedelta(minutes=thestrat.analyzers.TradeAnalyzer.get_analysis().len.average)
    ))
    print('Expectancy[%]         :{0:.2f}'.format(  # 期待値=平均利益×勝率＋平均損失×敗率
        thestrat.analyzers.TradeAnalyzer.get_analysis().won.pnl.average * win_rate
        + thestrat.analyzers.TradeAnalyzer.get_analysis().lost.pnl.average * lost_rate
    ))
    print('SQN                   :{0:.2f}'.format(  # SQN システムの評価値
        thestrat.analyzers.SQN.get_analysis().sqn
    ))

    print('Final Cash Value: %.2f' % cerebro.broker.getcash())
    cerebro.plot(volume=False)


# 引数を取り扱います。
# noinspection PyTypeChecker
def parse_args(p_args=None):
    parser = argparse.ArgumentParser(
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
        description=(
            'Captain America Back Test'
        )
    )
    parser.add_argument('--instrument', required=True, default='',
                        help='Target formatted instrument ex)"USD_JPY"')
    parser.add_argument('--granularity', required=True, default='',
                        help="Target granularity")
    return parser.parse_args(p_args)


if __name__ == '__main__':
    run_start()
