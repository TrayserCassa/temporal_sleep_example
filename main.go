package sleep

import (
    "strconv"
    "time"

    "go.temporal.io/sdk/workflow"
)

func KeepVersion2(ctx workflow.Context) error {
    // Wait one day
    err := workflow.Sleep(ctx, time.Hour*24)
    if err != nil {
        return err
    }

    var signalValue string
    signalChan := workflow.GetSignalChannel(ctx, "keep_version")
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
        // Could be 0 to 7 - The numbers stands for the day when we want to remind again. If it is 0 we can delete it immediately
        channel.Receive(ctx, &signalValue)
    })

    // Activity which is asking via slack if the user wants to keep the current version
    // We do some cancel context here, because if the User is not responding, we simply delete the version

    selector.Select(ctx)

    // Signal was sent, so we check the value
    daysToWait, err := strconv.Atoi(signalValue)
    if err != nil {
        return workflow.NewContinueAsNewError(ctx, KeepVersion2)
    }

    if daysToWait != 0 {
        err := workflow.Sleep(ctx, time.Hour*24*time.Duration(daysToWait-1))
        if err != nil {
            return err
        }
        return workflow.NewContinueAsNewError(ctx, KeepVersion2)
    }

    // Delete the version because the signal value was 0

    return nil
}
